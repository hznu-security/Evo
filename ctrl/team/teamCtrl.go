/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/09 14:12
 * 描述     : 队伍管理以及选手端请求
 */

package team

import (
	"Evo/auth"
	"Evo/config"
	"Evo/db"
	"Evo/model"
	"Evo/starry"
	"Evo/util"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

// TeamLogin 队伍登陆
func TeamLogin(c *gin.Context) {
	type loginForm struct {
		Name string `binding:"required,max=50"`
		Pwd  string `binding:"required,max=30"`
	}
	var form loginForm
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println(err.Error())
		util.Error(c, "绑定错误", nil)
	}
	var team model.Team
	if err = db.DB.Where("name = ?", form.Name).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Fail(c, "队伍不存在", nil)
			return
		} else {
			util.Error(c, "服务端错误", nil)
			return
		}
	}
	if team.Pwd != form.Pwd {
		util.Fail(c, "密码错误", nil)
		return
	}
	token, err := auth.ReleaseToken(team.ID, auth.TEAM)
	if err != nil {
		log.Println(err.Error())
	}
	db.DB.Model(model.Team{}).Where("id = ?", team.ID).Update("token", token)
	util.Success(c, "登陆成功", gin.H{
		"token": token,
		"id":    team.ID,
	})
}

type flagFrom struct {
	Flag string `json:"flag" binding:"required,max=255"`
}

// SubmitFlag 提交flag   需要解决flag重复提交问题
func SubmitFlag(c *gin.Context) {
	var form flagFrom
	err := c.ShouldBind(&form)
	if err != nil {
		util.Fail(c, "提交失败", nil)
	}
	teamId, isExist := c.Get("teamId") // 获取鉴权中间件放进去的teamId
	if !isExist {
		util.Fail(c, "队伍不存在", nil)
		return
	}
	var flag model.Flag
	db.DB.Where("round = ? AND flag = ?", config.ROUND_NOW, form.Flag).Find(&flag)
	// flag不正确,返回
	if flag.ID == 0 {
		util.Success(c, "flag不正确", nil)
		return
	}

	if flag.TeamId == teamId || flag.Round != config.ROUND_NOW {
		util.Success(c,"flag不正确",nil)
		return
	}

	// flag正确,判断是否提交过了
	var attack model.Attack
	attacker := teamId.(uint)
	db.DB.Where("attacker = ? AND round = ? AND game_box_id = ?", attacker, config.ROUND_NOW, flag.GameBoxId).First(&attack)
	if attack.ID != 0 {
		util.Success(c, "重复提交", nil)
		return
	}

	// flag正确
	var box model.GameBox
	db.DB.Where("id = ?", flag.GameBoxId).First(&box)
	if box.IsAttacked { // box已经被攻击过了

	}
	// box 没有被攻击过
	box.IsAttacked = true
	box.Score -= float64(config.ATTACK_SCORE)
	db.DB.Save(&box) // 更新靶机状态

	// 给大屏发消息
	starry.SendAttack(attacker, flag.TeamId)

	attack.Attacker = attacker
	attack.Round = config.ROUND_NOW
	attack.GameBoxId = box.ID
	attack.ChallengeId = flag.ChallengeID
	attack.TeamID = flag.TeamId // 被攻击者

	db.DB.Create(&attack)
	util.Success(c, "success", nil)
}

type info struct {
	Team  model.Team `json:"team"`
	Round uint       `json:"round"`
}

// Info 获取信息
func Info(c *gin.Context) {
	teamId := c.Query("teamId")
	var teamInfo model.Team
	db.DB.Where("id = ?", teamId).Select([]string{"name", "logo", "score", "token"}).First(&teamInfo)
	res := info{
		Team:  teamInfo,
		Round: config.ROUND_NOW,
	}
	util.Success(c, "success", gin.H{
		"info": res,
	})
}

// GetNotification 获取公告
func GetNotification(c *gin.Context) {
	var notifications []model.Notification
	db.DB.Find(&notifications)
	util.Success(c, "success", gin.H{
		"notifications": notifications,
	})
}

type gameBox struct {
	model.GameBox
	Title          string  `json:"challengeName"`
	Desc           string  `json:"desc"`
	ChallengeScore float64 `json:"challengeScore"`
}

// GetGameBox 获取队伍的靶机信息
func GetGameBox(c *gin.Context) {
	teamId := c.Query("teamId")
	// 查找靶机信息
	boxes := make([]gameBox, 0)
	db.DB.Model(&model.GameBox{}).Select(`game_boxes.id,game_boxes.port,
	game_boxes.ssh_port,game_boxes.ssh_user,game_boxes.ssh_pwd,game_boxes.score,game_boxes.is_down,
	game_boxes.is_attacked,challenges.title,challenges.desc,challenges.score as challenge_score`).
		Joins("inner join challenges on challenges.id = game_boxes.challenge_id and challenges.visible = ?", true).
		Where("game_boxes.team_id = ?", teamId).Scan(&boxes)
	util.Success(c, "success", gin.H{
		"boxes": boxes,
	})
}
