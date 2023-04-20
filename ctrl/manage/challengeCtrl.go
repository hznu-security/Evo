/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:接收管理题目的请求
 */

package manage

import (
	"Evo/db"
	"Evo/model"
	"Evo/util"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// PostChallenge 上传Challenge
func PostChallenge(c *gin.Context) {
	var form model.Challenge
	err := c.ShouldBind(&form)
	if err != nil {
		util.Fail(c, "参数绑定错误", nil)
		return
	}
	if form.AutoRefresh && form.Command == "" {
		util.Fail(c, "缺少刷新flag的命令", nil)
		return
	}
	var challenge model.Challenge

	db.DB.Where("title = ?", form.Title).First(&challenge)
	if challenge.ID != 0 {
		util.Fail(c, "题目已存在", gin.H{
			"challenge": form,
		})
		return
	}

	challenge.Title = form.Title
	challenge.Desc = form.Desc
	challenge.Score = form.Score
	challenge.AutoRefresh = form.AutoRefresh
	challenge.Command = form.Command
	//challenge.Type = form.Type

	err = db.DB.Create(&challenge).Error
	if err != nil {
		log.Println(err)
		util.Error(c, "添加失败", nil)
		return
	}
	log.Println("添加题目:", challenge.Title)
	util.Success(c, "添加成功", gin.H{
		"challenge": challenge,
	})
}

// GetChallenge 获取所有challenge
func GetChallenge(c *gin.Context) {
	challenges := make([]model.Challenge, 0)
	db.DB.Find(&challenges)
	util.Success(c, "查找成功", gin.H{
		"challenges": challenges,
	})
	return
}

// DelChallenge 当题目有依赖他的靶机时，不能删除题目
func DelChallenge(c *gin.Context) {
	challengeId := c.Query("challengeId")
	id, err := strconv.Atoi(challengeId)
	if err != nil {
		util.Fail(c, "参数格式有误", gin.H{
			"param": challengeId,
		})
		return
	}
	var challenge model.Challenge
	db.DB.Where("id = ?", id).First(&challenge)
	if challenge.ID == 0 {
		util.Fail(c, "challenge 不存在", nil)
		return
	}

	var count int64
	db.DB.Model(&model.GameBox{}).Where("challenge_id = ?", id).Count(&count)
	if count != 0 {
		util.Fail(c, "删除失败，有依赖于题目的靶机", nil)
		return
	}
	db.DB.Delete(&challenge)
	log.Println("删除题目:", challenge.Title)
	util.Success(c, "删除成功", nil)
}

type putChallengeForm struct {
	ChallengeId uint    `binding:"required"`
	Title       string  `binding:"required,max=100"`
	Desc        string  `binding:"required,max=255"`
	Score       float64 `binding:"required"`
	AutoRefresh bool
	Command     string `binding:"required,max=255"`
}

func PutChallenge(c *gin.Context) {
	var form putChallengeForm
	err := c.ShouldBind(&form)
	log.Println(form)
	if err != nil {
		log.Println(err.Error())
		util.Fail(c, "参数绑定错误", nil)
		return
	}

	if form.AutoRefresh && form.Command == "" {
		util.Fail(c, "缺少刷新flag的命令", nil)
		return
	}

	var challenge model.Challenge
	db.DB.Where("id = ?", form.ChallengeId).First(&challenge)
	if challenge.ID == 0 {
		util.Fail(c, "challenge不存在", nil)
		return
	}

	challenge.Title = form.Title
	challenge.Score = form.Score
	challenge.Command = form.Command
	challenge.AutoRefresh = form.AutoRefresh
	challenge.Desc = form.Desc
	//challenge.Type = form.Type

	db.DB.Save(&challenge)
	log.Println("修改题目:", challenge.Title)
	util.Success(c, "成功", gin.H{
		"challenge": challenge,
	})
}

func Visible(c *gin.Context) {
	challengeId := c.Query("challengeId")
	id, err := strconv.Atoi(challengeId)
	if err != nil {
		util.Fail(c, "参数错误", gin.H{
			"challengeId": id,
		})
		return
	}
	var challenge model.Challenge

	db.DB.Model(&model.Challenge{}).Where("id = ?", id).First(&challenge)
	if challenge.ID == 0 {
		util.Fail(c, "队伍不存在", nil)
	}
	err = db.DB.Model(&model.Challenge{}).Where("id = ?", id).
		Update("visible", !challenge.Visible).Error

	if err != nil {
		util.Error(c, "设置失败", nil)
		return
	}
	util.Success(c, "设置成功", nil)
	return
}
