/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/09 14:12
 * 描述     : 队伍管理以及选手端请求
 */

package ctrl

import (
	"Evo/auth"
	"Evo/config"
	"Evo/db"
	"Evo/model"
	"Evo/starry"
	"Evo/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"sort"
	"strconv"
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
	flag string `binding:"required,max=255"`
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
	db.DB.Where("round = ? AND flag = ", config.ROUND_NOW, form.flag).First(&flag)
	// flag不正确,返回
	if flag.ID == 0 {
		util.Success(c, "flag不正确", nil)
		return
	}

	if flag.TeamId == teamId || flag.Round != config.ROUND_NOW {

	}

	// flag正确,判断是否提交过了
	var attack model.Attack
	attacker := teamId.(uint)
	db.DB.Where("attacker = ? AND round = ? AND box_id = ?", attacker, config.ROUND_NOW, flag.GameBoxId).First(&attack)
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
	attack.BoxId = box.ID
	attack.ChallengeId = flag.ChallengeID
	attack.TeamID = flag.TeamId // 被攻击者

	db.DB.Create(&attack)
}

type info struct {
	Team  model.Team  `json:"team"`
	Round uint        `json:"round"`
	Token string      `json:"token"`
	Boxes []model.Box `json:"boxes"`
}

// Info 获取信息
func Info(c *gin.Context) {
	teamId := c.Query("teamId")
	var teamInfo model.Team
	db.DB.Where("id = ?", teamId).Select([]string{"name", "logo", "score", "token"}).First(&teamInfo)
	var box []model.Box
	db.DB.Where("team_id = ?", teamId).Select([]string{"port", "ssh_user", "ssh_pwd", "score", "is_down", "is_attacked"}).
		Find(&box)
	res := info{
		Team:  teamInfo,
		Round: config.ROUND_NOW,
		Token: teamInfo.Token,
		Boxes: box,
	}
	util.Success(c, "success", gin.H{
		"info": res,
	})
}

// GetRank 获取排名
func GetRank(c *gin.Context) {
	var team []model.Team
	db.DB.Select([]string{"id", "name", "logo", "score"}).Find(&team)
	// 根据分数排序后返回
	sort.Slice(team, func(i, j int) bool {
		return team[i].Score > team[j].Score
	})
	util.Success(c, "success", gin.H{
		"rank": team,
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

// PostTeam 添加队伍
func PostTeam(c *gin.Context) {
	var team model.Team
	// 绑定参数
	if err := c.ShouldBind(&team); err != nil {
		util.Fail(c, "绑定错误", nil)
	}

	if team.Name == "" {
		util.Fail(c, "参数错误", nil)
		return
	}

	db.DB.Where("name = ?", team.Name).First(&team)
	if team.ID != 0 {
		util.Fail(c, "队伍已存在", nil)
		return
	}

	//给队伍随机生成密码
	team.Pwd = auth.NewPwd()
	err := db.DB.Create(&team).Error
	if err != nil {
		util.Fail(c, "添加失败", nil)
		log.Println(err.Error())
	} else {
		util.Success(c, "添加成功", gin.H{
			"team": team,
		})
		log.Println("Add Team", team.Name)
	}
}

// PutTeam 修改队伍信息
func PutTeam(c *gin.Context) {

	type Form struct {
		TeamId uint   `json:"teamId" binding:"required"`
		Name   string `binding:"required,max=200"`
		Logo   string `binding:"max=255"`
	}

	var form Form
	err := c.ShouldBind(&form)
	if err != nil {
		util.Fail(c, "参数错误", nil)
		return
	}

	var team model.Team
	db.DB.Where("id = ?", form.TeamId).First(&team)
	if team.ID == 0 {
		util.Fail(c, "队伍不存在", nil)
		return
	}

	team.Name = form.Name
	team.Logo = form.Logo
	if err := db.DB.Save(&team).Error; err != nil {
		log.Println(err.Error())
		util.Fail(c, "保存失败", nil)
		return
	}
	util.Success(c, "修改成功", gin.H{
		"team": team,
	})
}

// GetTeam 列出所有队伍
func GetTeam(c *gin.Context) {
	teams := make([]model.Team, 0)
	db.DB.Find(&teams) //这里采用软删除，gorm自动忽视软删除过的内容
	util.Success(c, "查询成功", gin.H{
		"teams": teams,
	})
}

// DelTeam 删除队伍
func DelTeam(c *gin.Context) {
	teamIdStr := c.Query("teamId")
	teamId, err := strconv.Atoi(teamIdStr)
	if err != nil {
		util.Error(c, "参数错误", gin.H{
			"teamId": teamId,
		})
		log.Println(err.Error())
		return
	}
	var team model.Team

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&team, teamId).Error; err != nil {
			return err
		}
		if err = tx.Delete(&team).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Fail(c, "队伍不存在", gin.H{
				"teamId": teamId,
			})
			return
		} else {
			log.Println(err)
			util.Error(c, "删除失败", nil)
		}
		return
	}
	log.Printf("删除队伍 %s", team.Name)
	util.Success(c, "删除成功", nil)
}

// ResetPwd 重置队伍密码
func ResetPwd(c *gin.Context) {
	// 获得teamId
	teamIdStr := c.Query("teamId")
	if teamIdStr == "" {
		util.Fail(c, "参数错误", nil)
		return
	}
	teamId, err := strconv.Atoi(teamIdStr)
	if err != nil {
		util.Error(c, "服务端错误", nil)
		log.Println(err.Error())
		return
	}

	// 修改密码
	var team model.Team
	db.DB.Where("id = ?", teamId).First(&team)
	if team.ID == 0 {
		util.Fail(c, "队伍不存在", nil)
		return
	}
	team.Pwd = auth.NewPwd()
	db.DB.Save(&team)
	log.Printf("队伍 %s 重置密码", team.Name)
	util.Success(c, "重置成功", gin.H{
		"pwd": team.Pwd,
	})
}

// UploadLogo 上传队伍logo
func UploadLogo(c *gin.Context) {

	logo, err := c.FormFile("logo")
	if err != nil {
		util.Fail(c, "上传失败", nil)
		log.Println(err)
		return
	}
	ext := filepath.Ext(logo.Filename)
	if ext != ".png" && ext != ".jpg" {
		util.Fail(c, "图片格式不正确", nil)
		return
	}

	logoPath := viper.GetString("logo.path")
	err = util.TestAndMake(logoPath)
	if err != nil {
		util.Error(c, "上传失败", nil)
		log.Println(err)
		return
	}
	dst := util.GetRandomStr(10, logoPath, ext)
	err = c.SaveUploadedFile(logo, dst)
	if err != nil {
		util.Error(c, "上传失败", nil)
		log.Println(err)
		return
	}
	dst = dst[1:] // 去除路径开头的  .
	// 如果teamId 为空，表示是先上传logo再添加队伍的情况，如果teamId不为空，则是添加队伍，再上传logo
	teamId := c.PostForm("teamId")
	if teamId != "" {
		db.DB.Model(model.Team{}).Where("id = ?", teamId).Update("logo", dst)
	}
	util.Success(c, "success", gin.H{
		"path": dst,
	})
}
