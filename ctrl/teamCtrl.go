/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/09 14:12
 * 描述     : 队伍管理以及选手端请求
 */

package ctrl

import (
	"Evo/auth"
	"Evo/db"
	"Evo/model"
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		Error(c, "绑定错误", nil)
	}
	var team model.Team
	if err = db.DB.Where("name = ?", form.Name).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(c, "队伍不存在", nil)
			return
		} else {
			Error(c, "服务端错误", nil)
			return
		}
	}
	if team.Pwd != form.Pwd {
		Fail(c, "密码错误", nil)
		return
	}
	token, err := auth.ReleaseToken(team.ID, auth.TEAM)
	if err != nil {
		log.Println(err.Error())
	}
	Success(c, "登陆成功", gin.H{
		"token": token,
	})
}

// SubmitFlag 提交flag
func SubmitFlag(c *gin.Context) {

}

// Info 获取信息
func Info(c *gin.Context) {

}

// GetRank 获取排名
func GetRank(c *gin.Context) {

}

// GetNotification 获取公告
func GetNotification(c *gin.Context) {

}

// PostTeam 添加队伍
func PostTeam(c *gin.Context) {
	var team model.Team
	// 绑定参数
	if err := c.ShouldBind(&team); err != nil {
		Fail(c, "绑定错误", nil)
	}

	if team.Name == "" {
		Fail(c, "参数错误", nil)
		return
	}

	db.DB.Where("name = ?", team.Name).First(&team)
	if team.ID != 0 {
		Fail(c, "队伍已存在", nil)
		return
	}

	//给队伍随机生成密码
	team.Pwd = auth.NewPwd()
	err := db.DB.Create(&team).Error
	if err != nil {
		Fail(c, "添加失败", nil)
		log.Println(err.Error())
	} else {
		Success(c, "添加成功", nil)
		log.Println("Add Team", team.Name)
	}
}

// PutTeam 修改队伍信息
func PutTeam(c *gin.Context) {

	type Form struct {
		TeamId uint   `json:"teamId" binding:"required"`
		Name   string `binding:"required"`
		Logo   string
	}

	var form Form
	err := c.ShouldBind(&form)
	if err != nil {
		Fail(c, "参数错误", nil)
		return
	}

	var team model.Team
	db.DB.Where("id = ?", form.TeamId).First(&team)
	if team.ID == 0 {
		Fail(c, "队伍不存在", nil)
		return
	}

	team.Name = form.Name
	team.Logo = form.Logo
	if err := db.DB.Save(&team).Error; err != nil {
		log.Println(err.Error())
		Fail(c, "保存失败", nil)
		return
	}
	Success(c, "修改成功", nil)
}

// GetTeam 列出所有队伍
func GetTeam(c *gin.Context) {
	teams := make([]model.Team, 0)
	db.DB.Find(&teams) //这里采用软删除，gorm自动忽视软删除过的内容
	Success(c, "查询成功", gin.H{
		"teams": teams,
	})
}

// DelTeam 删除队伍
func DelTeam(c *gin.Context) {
	teamIdStr := c.Query("teamId")
	teamId, err := strconv.Atoi(teamIdStr)
	if err != nil {
		Error(c, "参数错误", gin.H{
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
			Fail(c, "队伍不存在", gin.H{
				"teamId": teamId,
			})
			return
		} else {
			log.Println(err)
			Error(c, "删除失败", nil)
		}
		return
	}
	Success(c, "删除成功", nil)
}

// ResetPwd 重置队伍密码
func ResetPwd(c *gin.Context) {
	// 获得teamId
	teamIdStr := c.Query("teamId")
	if teamIdStr == "" {
		Fail(c, "参数错误", nil)
		return
	}
	teamId, err := strconv.Atoi(teamIdStr)
	if err != nil {
		Error(c, "服务端错误", nil)
		log.Println(err.Error())
		return
	}

	// 修改密码
	var team model.Team
	db.DB.Where("id = ?", teamId).First(&team)
	if team.ID == 0 {
		Fail(c, "队伍不存在", nil)
		return
	}
	team.Pwd = auth.NewPwd()
	db.DB.Save(&team)
	Success(c, "重置成功", gin.H{
		"pwd": team.Pwd,
	})
}

// UploadLogo 上传队伍logo
func UploadLogo(c *gin.Context) {

}
