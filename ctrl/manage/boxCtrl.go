/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/12 11:19
 * 描述     ：管理靶机相关
 */

package manage

import (
	"Evo/auth"
	"Evo/db"
	"Evo/model"
	"Evo/service/box"
	"Evo/service/docker"
	"Evo/util"
	"encoding/json"
	"errors"
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

	ChallengePort map[uint]string `json:"challengePort" binding:"required"` // 题目id和题目的port
	TeamId        uint            `binding:"required"`
	Image         string          `binding:"required,max=100"` // 选择容器的名字
	SshPort       string          `json:"sshPort" binding:"required"`
	SshUser       string          `json:"sshUser" binding:"required"`
	SshPwd        string          `json:"sshPwd"`
}

// TODO 加检查，别重复部署
func PostBox(c *gin.Context) {
	// 绑定参数
	var boxForm BoxForm
	if err := c.ShouldBind(&boxForm); err != nil {
		log.Println(err)
		util.Fail(c, "参数绑定失败", nil)
		return
	}

	// 检查题目是否存在，靶机是否已存在，port格式是否正确
	if err := checkChallengePort(boxForm.ChallengePort, boxForm.TeamId); err != nil {
		util.Fail(c, err.Error(), nil)
		return
	}

	err := docker.CheckImage(boxForm.Image)
	if err != nil {
		if err.Error() == "镜像不存在" {
			util.Fail(c, "镜像不存在", nil)
		} else {
			log.Println(err.Error())
			util.Error(c, "检查镜像失败", nil)
		}
		return
	}

	var newBox model.Box
	var challenge model.Challenge
	portMap := make(nat.PortMap)

	cNameBuilder := strings.Builder{}

	var gameBoxes []model.GameBox

	for challengeId, port := range boxForm.ChallengePort {
		// 解析port   port形式:8080:8080,9090:9090
		pMap := docker.ParsePort(port)

		gameBoxPort := ""
		db.DB.Where("id = ?", challengeId).First(&challenge)
		for k, v := range pMap {
			portMap[k] = v
			gameBoxPort += string(k)
			gameBoxPort = gameBoxPort + ","
		}
		cNameBuilder.WriteString("C")
		cNameBuilder.WriteString(strconv.Itoa(int(challengeId)))

		// 新增一个靶机
		// 靶机名字市题目名+队伍名
		gameBoxes = append(gameBoxes, model.GameBox{
			TeamId:      boxForm.TeamId,
			SshPort:     boxForm.SshPort,
			SshUser:     boxForm.SshUser,
			Port:        gameBoxPort,
			Score:       challenge.Score,
			ChallengeID: challengeId,
		})
	}
	cNameBuilder.WriteString("T")
	cNameBuilder.WriteString(strconv.Itoa(int(boxForm.TeamId)))

	cName := cNameBuilder.String()

	// 开始启动容器
	// 传入镜像名，容器名，端口映射
	sshPortMap := docker.ParsePort(boxForm.SshPort + ":22")
	for k, v := range sshPortMap {
		portMap[k] = v
	}
	err = docker.StartContainer(boxForm.Image, cName, &portMap)
	if err != nil {
		log.Println(err)
		errr := docker.RemoveContainer(cName)
		log.Println(errr)
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
	newBox.SshPort = boxForm.SshPort
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
		gameBoxes[i].SshPwd = newBox.SshPwd
	}

	// 插入gamebox
	db.DB.Create(&gameBoxes)

	util.Success(c, "success", nil)
}

type gameBox struct {
	ID       uint
	CName    string  `json:"CName"`
	Port     string  `json:"port"`
	SshPort  string  `json:"sshPort"`
	SshUser  string  `json:"sshUser"`
	SshPwd   string  `json:"sshPwd"`
	Score    float64 `json:"score"`
	IsDown   bool    `json:"isDown"`
	IsAttack bool    `json:"isAttack"`
	Name     string  `json:"TeamName"`      // 队伍名
	Title    string  `json:"challengeName"` // 题目名
}

func GetBox(c *gin.Context) {
	boxes := make([]gameBox, 0)
	db.DB.Model(&model.GameBox{}).Select(`game_boxes.id,game_boxes.c_name,game_boxes.port,
	game_boxes.ssh_port,game_boxes.ssh_user,game_boxes.ssh_pwd,game_boxes.score,game_boxes.is_down,
	game_boxes.is_attacked,challenges.title,teams.name`).
		Joins("inner join challenges on challenges.id = game_boxes.challenge_id").
		Joins("inner join teams on teams.id = game_boxes.team_id").Scan(&boxes)
	util.Success(c, "success", gin.H{
		"boxes": boxes,
	})
}

type putBoxForm struct {
	BoxId       uint `binding:"required"`
	ChallengeId uint `binding:"required"`
	TeamId      uint `binding:"required"`
}

// 靶机创建后直接不允许修改，直接删容器好了
func PutBox(c *gin.Context) {
	//var form putBoxForm
	//err := c.ShouldBind(&form)
	//if err != nil {
	//	log.Println(err)
	//	util.Fail(c, "参数绑定失败", nil)
	//	return
	//}
	//var box model.Box
	//db.DB.Where("id = ?", form.BoxId).First(&box)
	//if box.ID == 0 {
	//	util.Fail(c, "靶机不存在", nil)
	//	return
	//}
	//box.ChallengeID = form.ChallengeId
	//box.TeamId = form.TeamId
	//db.DB.Save(&box)
	//
	//log.Println("修改靶机:", box.Name)
	//util.Success(c, "success", gin.H{
	//	"box": box,
	//})
}

// DelBox 删除靶机，不支持修改靶机信息，想修改直接调用到这里删除容器
func DelBox(c *gin.Context) {
	boxId := c.Query("gameBoxId")
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

	var container model.Box
	db.DB.Where("name = ?", gameBox.CName).First(&container)
	if container.ID == 0 {
		errMsg := fmt.Sprintf("容器:%s不存在", gameBox.CName)
		log.Println(errMsg)
		util.Error(c, errMsg, nil)
	}
	// 删除容器
	if err := docker.RemoveContainer(container.Name); err != nil {
		log.Println(err.Error())
		util.Error(c, fmt.Sprintf("删除容器:%s失败", container.Name), nil)
	}

	// 删除该容器信息
	db.DB.Delete(&container)
	// 删除所有与该容器相关的game_box
	db.DB.Delete(&model.GameBox{}, "c_name = ?", gameBox.CName)
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

// 给端口，用户明，密码。
func testSSH(cName, user, pwd string) (string, error) {
	var box model.Box
	db.DB.Where("name = ?", cName).First(&box)
	res, err := util.SSHExec(box.SshPort, user, pwd, "whoami")
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

func checkPortFormat(port string) error {
	pms := strings.Split(port, ",")
	for _, pm := range pms {
		if !strings.Contains(pm, ":") {
			return errors.New(fmt.Sprintf("port格式错误 %s", port))
		}
		m := strings.Split(pm, ":")
		if len(m) != 2 {
			return errors.New(fmt.Sprintf("port格式错误 %s", port))
		}
		// 判断是不是整数
		for _, p := range m {
			if _, err := strconv.Atoi(p); err != nil {
				return errors.New(fmt.Sprintf("port格式错误 %s", p))
			}
		}
	}
	return nil
}

func checkChallengePort(challengePort map[uint]string, teamId uint) error {
	var box model.GameBox
	var challenge model.Challenge
	for challengeId, port := range challengePort {
		db.DB.Where("challenge_id = ? and team_id = ?", challengeId, teamId).First(&box)
		if box.ID != 0 {
			return errors.New(fmt.Sprintf("靶机已存在,队伍:%v 题目:%v", teamId, challengeId))
		}
		db.DB.Where("id = ?", challengeId).First(&challenge)
		if challenge.ID == 0 {
			return errors.New(fmt.Sprintf("题目:%v不存在", challengeId))
		}
		err := checkPortFormat(port)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateScore(c *gin.Context) {
	updateScore()
}

// TODO  这里有点问题 按说比赛开始后不能更新分数了
func updateScore() {
	teams := make([]model.Team, 0)
	db.DB.Find(&teams)
	for i := 0; i < len(teams); i++ {
		gameBoxes := make([]model.GameBox, 0)
		db.DB.Select([]string{"score"}).Where("team_id = ?", teams[i].ID).Find(&gameBoxes)
		for j := 0; j < len(gameBoxes); j++ {
			teams[i].Score += gameBoxes[j].Score
		}
	}
	db.DB.Save(&teams)
}
