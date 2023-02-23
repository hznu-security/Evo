/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/12 11:19
 * 描述     ：管理靶机相关
 */

package ctrl

import (
	"Evo/auth"
	"Evo/db"
	"Evo/model"
	"Evo/service/box"
	"Evo/service/docker"
	"Evo/util"
	"encoding/json"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
)

// BoxForm 添加靶机
type BoxForm struct {
	//Port          string          `binding:"required,max=100"`
	//ChallengeId   uint            `binding:"required"`

	ChallengePort map[uint]string `binding:"required"` // 题目id和题目的port
	TeamId        uint            `binding:"required"`
	Image         string          `binding:"required,max=100"` // 选择容器的名字
	SshPort       string          `json:"sshPort" binding:"required"`
	SshUser       string          `json:"sshUser" binding:"required"`
	SshPwd        string          `json:"sshPwd" binding:"required"`
}

func PostBox(c *gin.Context) {
	// 绑定参数
	var boxForm BoxForm
	if err := c.ShouldBind(&boxForm); err != nil {
		log.Println(err)
		util.Fail(c, "参数绑定失败", nil)
		return
	}

	// 检查一遍
	var newBox model.Box
	var challenge model.Challenge
	portMap := make(nat.PortMap)

	cNameBuilder := strings.Builder{}

	var gameBoxes []model.GameBox

	for challengeId, port := range boxForm.ChallengePort {
		db.DB.Where("challenge_id = ? and team_id = ?", challengeId, boxForm.TeamId).First(&newBox)
		if newBox.ID != 0 {
			util.Fail(c, fmt.Sprintf("靶机已存在,队伍:%v 题目:%v", boxForm.TeamId, challengeId), nil)
			return
		}
		db.DB.Where("id = ?", challengeId).First(&challenge)
		if challenge.ID == 0 {
			util.Fail(c, fmt.Sprintf("题目:%v不存在", challengeId), nil)
			return
		}
		// 解析port   port形式:8080:8080,9090:9090
		pMap := docker.ParsePort(port)

		gameBoxPort := ""

		for k, v := range pMap {
			portMap[k] = v
			gameBoxPort += string(k)
			gameBoxPort = gameBoxPort + ","
		}
		cNameBuilder.WriteString("C")
		cNameBuilder.WriteRune(rune(challengeId))

		gameBoxes = append(gameBoxes, model.GameBox{
			TeamId:  boxForm.TeamId,
			Name:    challenge.Title,
			SshPort: boxForm.SshPort,
			SshUser: boxForm.SshUser,
			SshPwd:  boxForm.SshPwd,
			Port:    gameBoxPort,
			Score:   challenge.Score,
		})

	}
	cNameBuilder.WriteString("T")
	cNameBuilder.WriteRune(rune(boxForm.TeamId))

	cName := cNameBuilder.String()

	// 开始启动容器
	// 传入镜像名，容器名，端口映射
	err := docker.StartContainer(boxForm.Image, cName, &portMap)
	if err != nil {
		log.Println(err)
		err := docker.RemoveContainer(cName)
		log.Println(err)
		util.Error(c, "启动靶机失败", gin.H{
			"error": err.Error(),
		})
		return
	}

	port, _ := json.Marshal(portMap)
	sshPwd := auth.NewPwd()

	newBox.TeamId = boxForm.TeamId
	newBox.SshUser = boxForm.SshUser
	newBox.SshPwd = sshPwd
	newBox.Name = cName
	newBox.Port = string(port)
	err = docker.SetContainerSSH(cName, boxForm.SshUser, sshPwd)
	if err != nil {
		log.Println(err)
		err = docker.RemoveContainer(cName)
		if err != nil {
			log.Println(err)
		}
		util.Error(c, "设置ssh账号失败,请手动操作", nil)
		return
	}
	log.Println("容器", cName, "启动")

	err = db.DB.Create(&newBox).Error
	if err != nil {
		log.Println(err)
		err = docker.RemoveContainer(cName)
		if err != nil {
			log.Println(err)
		}
		util.Error(c, "启动成功，数据库异常", nil)
		return
	}

	// 开始写入gamebox
	for i := 0; i < len(gameBoxes); i++ {
		gameBoxes[i].CName = cName
	}

	util.Success(c, "success", nil)
}

func GetBox(c *gin.Context) {
	boxes := make([]model.GameBox, 0)
	db.DB.Find(&boxes)
	util.Success(c, "success", gin.H{
		"boxes": boxes,
	})
}

type putBoxForm struct {
	BoxId       uint `binding:"required"`
	ChallengeId uint `binding:"required"`
	TeamId      uint `binding:"required"`
}

// TODO
func PutBox(c *gin.Context) {
	var form putBoxForm
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println(err)
		util.Fail(c, "参数绑定失败", nil)
		return
	}
	var box model.Box
	db.DB.Where("id = ?", form.BoxId).First(&box)
	if box.ID == 0 {
		util.Fail(c, "靶机不存在", nil)
		return
	}
	box.ChallengeID = form.ChallengeId
	box.TeamId = form.TeamId
	db.DB.Save(&box)

	log.Println("修改靶机:", box.Name)
	util.Success(c, "success", gin.H{
		"box": box,
	})
}

// TODO 如果一个容器对应的最后一个gamebox被删除了，就删除这个容器
// DelBox 移除靶机,这里不删除容器
func DelBox(c *gin.Context) {
	boxId := c.Query("boxId")
	id, err := strconv.Atoi(boxId)
	if err != nil {
		util.Fail(c, "参数错误", gin.H{
			"param": boxId,
		})
		return
	}
	var gameBox model.GameBox
	db.DB.Where("id = ?", id).First(&gameBox)
	if gameBox.ID == 0 {
		util.Fail(c, "靶机不存在", nil)
		return
	}

	db.DB.Delete(&gameBox)
	log.Println("删除靶机:", gameBox.Name)
	util.Success(c, "删除成功", nil)
}

// ResetBox 重置所有靶机的状态，连带分数
func ResetBox(c *gin.Context) {
	err := box.ResetAllStatus()
	if err != nil {
		log.Println(err)
		util.Error(c, "更新失败", nil)
		return
	}

	err = box.ResetAllScore()
	if err != nil {
		log.Println(err)
		util.Error(c, "更新失败", nil)
		return
	}
	log.Println("重置靶机状态和分数")
	util.Success(c, "重置成功", nil)
}

// TestSSH 测试ssh链接
func TestSSH(c *gin.Context) {
	boxId := c.Query("boxId")
	id, err := strconv.Atoi(boxId)
	if err != nil {
		util.Fail(c, "参数格式有误", gin.H{
			"param": boxId,
		})
	}
	var box model.GameBox
	db.DB.Where("id = ?", id).First(&box)
	if box.ID == 0 {
		util.Fail(c, "id错误，靶机不存在", nil)
		return
	}

	res, err := testSSH(box.CName, box.SshUser, box.SshPwd)
	if err != nil {
		log.Println(err)
		util.Fail(c, "测试失败", gin.H{
			"res": "whoami:" + res,
		})
		return
	}
	log.Println("测试ssh")
	util.Success(c, "测试成功", gin.H{
		"res": "whoami:" + res,
	})
}

// TestSSHAll 测试所有靶机的ssh连接
func TestSSHAll(c *gin.Context) {
	faliedBoxes := make([]model.Box, 0)
	boxes := make([]model.Box, 0)
	db.DB.Select([]string{"ssh_user", "ssh_pwd", "name"}).Find(&boxes)

	for _, box := range boxes {
		if _, err := testSSH(box.Name, box.SshUser, box.SshPwd); err != nil {
			faliedBoxes = append(faliedBoxes, box)
		}
	}
	if len(faliedBoxes) == 0 {
		util.Success(c, "全部靶机ssh正常", nil)
		return
	} else {
		util.Fail(c, "存在靶机未通过", gin.H{
			"boxes": faliedBoxes,
		})
		return
	}
}

func testSSH(name, user, pwd string) (string, error) {
	ip, err := docker.GetIp(name)
	if err != nil {
		return "", err
	}
	res, err := util.SSHExec(ip, user, pwd, "whoami")
	if !strings.Contains(res, user) {
		return res, fmt.Errorf("测试失败")
	}
	if err != nil {
		return "", err
	}
	return res, err
}

func ReFreshFlag(c *gin.Context) {

}
