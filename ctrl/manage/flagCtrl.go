/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:管理flag相关
 */

package manage

import (
	"Evo/db"
	"Evo/model"
	"Evo/service/game"
	"Evo/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

// PostFlag 上传flag,针对非awd题目
func PostFlag(c *gin.Context) {
	type flagFrom struct {
		Flag        string `binding:"required,max=255"`
		ChallengeID uint   `binding:"required"`
	}
	var form flagFrom
	err := c.ShouldBind(&form)
	if err != nil {
		util.Fail(c, "参数有误", nil)
		return
	}
	flag := model.Flag{
		Flag:        form.Flag,
		ChallengeID: form.ChallengeID,
	}
	db.DB.Save(&flag)
	util.Success(c, "success", nil)
}

// 清空falg表
func clearFlag() {
	sql := fmt.Sprintf("TRUNCATE TABLE %s;", "flags")
	log.Println(sql)
	db.DB.Exec(sql)
}

// GenerateFlag 生成flag,针对awd题目
func GenerateFlag(c *gin.Context) {
	// 先清空flag
	clearFlag()
	err := game.GenerateFlag()
	if err != nil {
		util.Error(c, "生成flag失败", nil)
		return
	}
	util.Success(c, "生成flag成功", nil)
}

// ExportFlag 导出flag
func ExportFlag(c *gin.Context) {

}

type FlagPage struct {
	PageNum     int `json:"pageNum" binding:"required"`
	PageSize    int `json:"pageSize" binding:"required"`
	TeamId      int `json:"teamId"`
	ChallengeId int `json:"ChallengeId"`
	Round       int `json:"Round"`
}

func GetFlag(c *gin.Context) {
	var form FlagPage
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println(err.Error())
		util.Fail(c, "参数绑定失败", nil)
		return
	}
	m := make(map[string]interface{})
	if form.TeamId != 0 {
		m["team_id"] = form.TeamId
	}
	if form.ChallengeId != 0 {
		m["challenge_id"] = form.ChallengeId
	}
	if form.Round != 0 {
		m["Round"] = form.Round
	}
	var flags []model.Flag
	var count int64
	err = db.DB.Where(m).Limit(form.PageSize).Offset((form.PageNum - 1) * form.PageSize).Find(&flags).
		Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		log.Println(err.Error())
		util.Error(c, "查询错误", nil)
		return
	}
	util.Success(c, "success", gin.H{
		"flags": flags,
		"count": count,
	})
}
