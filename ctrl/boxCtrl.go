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
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

// BoxForm 添加靶机
type BoxForm struct {
	ChallengeId uint   `binding:"required"`
	TeamId      uint   `binding:"required"`
	Ip          string `binding:"required,max=30"`
	Image       string `binding:"required,max=30"`
	Port        string `binding:"required"`
	SshUser     string `json:"sshUser" binding:"required"`
}

func PostBox(c *gin.Context) {
	// 绑定参数
	var boxForm BoxForm
	if err := c.ShouldBind(&boxForm); err != nil {
		log.Println(err)
		Fail(c, "参数绑定失败", nil)
		return
	}

	var box model.Box
	db.DB.Where("challenge_id = ? and team_id = ?", boxForm.ChallengeId, boxForm.TeamId).First(&box)
	if box.ID != 0 {
		Fail(c, "靶机已存在", nil)
		return
	}

	var challenge model.Challenge
	db.DB.Where("id = ?", boxForm.ChallengeId).First(&challenge)
	if challenge.ID == 0 {
		Fail(c, "队伍不存在", nil)
		return
	}
	// 容器名为题目名+队伍id
	name := challenge.Title + strconv.Itoa(int(boxForm.TeamId))
	// 获取容器所在网络名
	network := viper.GetString("docker.network")
	// 传入镜像名，ip，网络名
	err := docker.StartContainer(boxForm.Image, name, network, boxForm.Ip)
	if err != nil {
		log.Println(err)
		Error(c, "启动靶机失败", nil)
		return
	}

	sshPwd := auth.NewPwd()

	box.ChallengeID = boxForm.ChallengeId
	box.TeamId = boxForm.TeamId
	box.Ip = boxForm.Ip
	box.Port = boxForm.Port
	box.Score = challenge.Score
	box.SshUser = boxForm.SshUser
	box.SshPwd = sshPwd
	box.Name = name

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

	err = db.DB.Create(&box).Error
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

	err = docker.RemoveContainer(box.Name)
	if err != nil {
		log.Println(err)
		Error(c, "移除失败", nil)
		return
	}
	db.DB.Delete(&box)
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

	res, err := util.SSHExec(box.Ip, box.SshUser, box.SshPwd, "whoami")
	if err != nil {
		log.Println(err)
		Fail(c, "测试失败", nil)
		return
	}
	if res != box.SshUser {
		log.Printf("靶机：%s:whoami用户不正确", box.ID)
		Fail(c, "whoami用户不正确", nil)
	}
	log.Println("测试ssh")
	Success(c, "测试成功", gin.H{
		"res": "whoami:" + res,
	})
}

func ReFreshFlag(c *gin.Context) {

}
