package manage

import (
	"Evo/config"
	"Evo/cron"
	"Evo/db"
	"Evo/model"
	"Evo/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func Check(c *gin.Context) {

	if !cron.IsOn() {
		util.Success(c, "比赛未在进行", nil)
		return
	}

	type checkForm struct {
		GameBoxId uint `binding:"required" json:"GameBoxId"`
	}
	var form checkForm
	if err := c.ShouldBind(&form); err != nil {
		util.Fail(c, "参数绑定失败", nil)
		return
	}

	// 记入down表
	var down model.Down
	db.DB.Model(&model.Down{}).Where(&model.Down{
		GameBoxId: form.GameBoxId,
		Round:     config.ROUND_NOW,
	}).First(&down)
	if down.ID != 0 {
		util.Fail(c, "check repeat", nil)
		return
	}

	var gBox model.GameBox
	db.DB.Model(&model.GameBox{}).Where("id = ?", form.GameBoxId).First(&gBox)
	if gBox.ID == 0 {
		util.Fail(c, "靶机不存在", nil)
		return
	}
	if !gBox.Visible {
		util.Fail(c, "靶机不可见", nil)
		return
	}
	db.DB.Model(&model.GameBox{}).Where("id = ?", gBox.ID).Update("is_down", true)
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model.Down{
			TeamId:      gBox.TeamId,
			ChallengeId: gBox.ChallengeID,
			Round:       config.ROUND_NOW,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf(err.Error())
		util.Error(c, "记录宕机失败", nil)
		return
	}
	util.Success(c, "success", nil)
	var team model.Team
	var challenge model.Challenge
	db.DB.Where("id = ?", gBox.TeamId).First(&team)
	db.DB.Where("id = ?", gBox.ChallengeID).First(&challenge)
	log.Printf("check down Team:%s Challenge:%s", team.Name, challenge.Title)
}
