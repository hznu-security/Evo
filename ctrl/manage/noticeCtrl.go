/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		：处理通知相关的请求
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

// PostNotice 上传通知
func PostNotice(c *gin.Context) {
	var notice model.Notification
	err := c.ShouldBind(&notice)
	if err != nil {
		util.Fail(c, "参数错误", gin.H{
			"notice": notice,
		})
		return
	}

	err = db.DB.Create(&notice).Error
	if err != nil {
		log.Println(err)
		util.Error(c, "插入失败", nil)
		return
	}

	util.Success(c, "success", gin.H{
		"notification": notice,
	})
}

type putNoticeForm struct {
	NotificationId uint   `binding:"required"`
	Title          string `binding:"required,max=100"`
	Content        string `binding:"required,max=255"`
}

// PutNotice 修改通知
func PutNotice(c *gin.Context) {
	var form putNoticeForm
	err := c.ShouldBind(&form)
	if err != nil {
		util.Fail(c, "参数绑定错误", gin.H{
			"notification": form,
		})
		return
	}

	var notice model.Notification
	db.DB.Where("id = ?", form.NotificationId).First(&notice)
	if notice.ID == 0 {
		util.Fail(c, "通知不存在", gin.H{
			"notificationId": form.NotificationId,
		})
		return
	}
	notice.Title = form.Title
	notice.Content = form.Content
	db.DB.Save(&notice)
	util.Success(c, "修改成功", gin.H{
		"notification": notice,
	})
}

// GetNotice 获取通知
func GetNotice(c *gin.Context) {
	notices := make([]model.Notification, 0)
	db.DB.Find(&notices)
	util.Success(c, "success", gin.H{
		"notifications": notices,
	})
}

// DelNotice 删除通知
func DelNotice(c *gin.Context) {
	noticeId := c.Query("notificationId")
	id, err := strconv.Atoi(noticeId)
	if err != nil {
		util.Fail(c, "参数错误", gin.H{
			"noticeId": noticeId,
		})
		return
	}
	var notice model.Notification
	db.DB.First(&notice, id)
	if notice.ID == 0 {
		util.Fail(c, "通知不存在", gin.H{
			"notificationId": id,
		})
		return
	}
	db.DB.Delete(&notice)
	util.Success(c, "success", nil)
}
