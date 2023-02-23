package manage

import (
	"Evo/auth"
	"Evo/db"
	"Evo/model"
	"Evo/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"strconv"
)

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

	if form.Name != "" {
		team.Name = form.Name

	}
	if form.Logo != "" {
		team.Logo = form.Logo
	}

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
