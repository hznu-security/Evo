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
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
)

// BoxForm 添加靶机
type BoxForm struct {
	ChallengeId uint `binding:"required"`
	TeamId      uint `binding:"required"`
	//Ip          string `binding:"required,max=30"`
	Image   string `binding:"required,max=30"`
	Port    string `binding:"required,max=100"`
	SshUser string `json:"sshUser" binding:"required"`
}

func PostBox(c *gin.Context) {
	// 绑定参数
	var boxForm BoxForm
	if err := c.ShouldBind(&boxForm); err != nil {
		log.Println(err)
		Fail(c, "参数绑定失败", nil)
		return
	}

	var newBox model.Box
	db.DB.Where("challenge_id = ? and team_id = ?", boxForm.ChallengeId, boxForm.TeamId).First(&newBox)
	if newBox.ID != 0 {
		Fail(c, "靶机已存在", nil)
		return
	}

	var challenge model.Challenge
	db.DB.Where("id = ?", boxForm.ChallengeId).First(&challenge)
	if challenge.ID == 0 {
		Fail(c, "队伍不存在", nil)
		return
	}

	portMap := docker.ParsePort(boxForm.Port)

	// 容器名为题目名+队伍id
	name := challenge.Title + strconv.Itoa(int(boxForm.TeamId))
	// 传入镜像名，ip，网络名
	err := docker.StartContainer(boxForm.Image, name, &portMap)
	if err != nil {
		log.Println(err)
		errr := docker.RemoveContainer(name)
		log.Println(errr)
		Error(c, "启动靶机失败", gin.H{
			"error": err.Error(),
		})
		return
	}

	port, _ := json.Marshal(portMap)

	sshPwd := auth.NewPwd()

	newBox.ChallengeID = boxForm.ChallengeId
	newBox.TeamId = boxForm.TeamId
	//newBox.Ip = boxForm.Ip
	newBox.Score = challenge.Score
	newBox.SshUser = boxForm.SshUser
	newBox.SshPwd = sshPwd
	newBox.Name = name
	newBox.Port = string(port)
	err = docker.SetContainerSSH(name, boxForm.SshUser, sshPwd)
	if err != nil {
		log.Println(err)
		err = docker.RemoveContainer(name)
		if err != nil {
			log.Println(err)
		}
		Error(c, "设置ssh账号失败,请手动操作", nil)
		return
	}
	log.Println("靶机", name, "启动")

	err = db.DB.Create(&newBox).Error
	if err != nil {
		log.Println(err)
		err = docker.RemoveContainer(name)
		if err != nil {
			log.Println(err)
		}
		Error(c, "启动成功，数据库异常", nil)
		return
	}
	Success(c, "success", nil)
}

func GetBox(c *gin.Context) {
	boxes := make([]model.Box, 0)
	db.DB.Find(&boxes)
	Success(c, "success", gin.H{
		"boxes": boxes,
	})
}

type putBoxForm struct {
	BoxId       uint `binding:"required"`
	ChallengeId uint `binding:"required"`
	TeamId      uint `binding:"required"`
}

// PutBox 修改靶机信息
func PutBox(c *gin.Context) {
	var form putBoxForm
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println(err)
		Fail(c, "参数绑定失败", nil)
		return
	}
	var box model.Box
	db.DB.Where("id = ?", form.BoxId).First(&box)
	if box.ID == 0 {
		Fail(c, "靶机不存在", nil)
		return
	}
	box.ChallengeID = form.ChallengeId
	box.TeamId = form.TeamId
	db.DB.Save(&box)

	log.Println("修改靶机:", box.Name)
	Success(c, "success", gin.H{
		"box": box,
	})
}

// DelBox 移除靶机
func DelBox(c *gin.Context) {
	boxId := c.Query("boxId")
	id, err := strconv.Atoi(boxId)
	if err != nil {
		Fail(c, "参数错误", gin.H{
			"param": boxId,
		})
		return
	}
	var box model.Box
	db.DB.Where("id = ?", id).First(&box)
	if box.ID == 0 {
		Fail(c, "靶机不存在", nil)
		return
	}

	db.DB.Delete(&box)
	err = docker.RemoveContainer(box.Name)
	if err != nil {
		log.Println(err)
		Error(c, err.Error(), nil)
		return
	}
	log.Println("删除靶机:", box.Name)
	Success(c, "删除成功", nil)
}

// ResetBox 重置所有靶机的状态，连带分数
func ResetBox(c *gin.Context) {
	err := box.ResetAllStatus()
	if err != nil {
		log.Println(err)
		Error(c, "更新失败", nil)
		return
	}

	err = box.ResetAllScore()
	if err != nil {
		log.Println(err)
		Error(c, "更新失败", nil)
		return
	}
	log.Println("重置靶机状态和分数")
	Success(c, "重置成功", nil)
}

// TestSSH 测试ssh链接
func TestSSH(c *gin.Context) {
	boxId := c.Query("boxId")
	id, err := strconv.Atoi(boxId)
	if err != nil {
		Fail(c, "参数格式有误", gin.H{
			"param": boxId,
		})
	}
	var box model.Box
	db.DB.Where("id = ?", id).First(&box)
	if box.ID == 0 {
		Fail(c, "id错误，靶机不存在", nil)
		return
	}

	res, err := testSSH(box.Name, box.SshUser, box.SshPwd)
	if err != nil {
		log.Println(err)
		Fail(c, "测试失败", gin.H{
			"res": "whoami:" + res,
		})
		return
	}
	log.Println("测试ssh")
	Success(c, "测试成功", gin.H{
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
		Success(c, "全部靶机ssh正常", nil)
		return
	} else {
		Fail(c, "存在靶机未通过", gin.H{
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
