/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:管理flag相关
 */

package ctrl

import (
	"Evo/db"
	"Evo/model"
	"Evo/service/game"
	"github.com/gin-gonic/gin"
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
		Fail(c, "参数有误", nil)
		return
	}
	flag := model.Flag{
		Flag:        form.Flag,
		ChallengeID: form.ChallengeID,
	}
	db.DB.Save(&flag)
	Success(c, "success", nil)
}

// GenerateFlag 生成flag,针对awd题目
func GenerateFlag(c *gin.Context) {
	err := game.GenerateFlag()
	if err != nil {
		Error(c, "生成flag失败", nil)
		return
	}
	Success(c, "生成flag成功", nil)
}

// ExportFlag 导出flag
func ExportFlag(c *gin.Context) {

}

func FilterFlag(c *gin.Context) {

}
