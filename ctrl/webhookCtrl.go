/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:处理webhook相关接口
 */

package ctrl

import (
	"Evo/db"
	"Evo/model"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func PostWebhook(c *gin.Context) {
	var form model.Webhook
	err := c.ShouldBind(&form)
	if err != nil {
		Fail(c, "参数绑定失败", nil)
		return
	}

	var webhook model.Webhook
	db.DB.Where("url = ?", form.Url).First(&webhook)
	if webhook.ID != 0 {
		Fail(c, "webhook已存在", gin.H{
			"url": form.Url,
		})
		return
	}

	webhook = form
	err = db.DB.Create(&webhook).Error
	if err != nil {
		log.Println(err)
		Error(c, "创建出错", nil)
		return
	}
	log.Println("添加webhook:", webhook)
	Success(c, "添加成功", gin.H{
		"webhook": webhook,
	})
}

type putWebhookForm struct {
	WebhookId uint    `binding:"required"`
	Url       string  `binding:"required,max=255"`
	Type      string  `binding:"required,max=30"`
	Time      uint    `binding:"required"`
	Timeout   float64 `binding:"required"`
}

func PutWebhook(c *gin.Context) {
	var form putWebhookForm
	err := c.ShouldBind(&form)
	if err != nil {
		Fail(c, "参数绑定有误", gin.H{
			"webhook": form,
		})
		return
	}

	var webhook model.Webhook
	db.DB.Where("id = ?", form.WebhookId).First(&webhook)
	if webhook.ID == 0 {
		Fail(c, "webhook不存在", gin.H{
			"webhook": form,
		})
		return
	}

	webhook.Url = form.Url
	webhook.Type = form.Type
	webhook.Time = form.Time
	webhook.Timeout = form.Timeout
	err = db.DB.Save(&webhook).Error
	if err != nil {
		log.Println(err)
		Error(c, "更新失败", gin.H{
			"webhook": form,
		})
		return
	}
	log.Println("修改webhook:", webhook)
	Success(c, "success", nil)
}

func GetWebhook(c *gin.Context) {
	webhooks := make([]model.Webhook, 0)
	db.DB.Find(&webhooks)
	Success(c, "success", gin.H{
		"webhooks": webhooks,
	})
}

func DelWebhook(c *gin.Context) {
	webhookId := c.Query("webhookId")
	id, err := strconv.Atoi(webhookId)
	if err != nil {
		Fail(c, "参数有误", gin.H{
			"webhookUId": webhookId,
		})
		return
	}
	var webhook model.Webhook
	db.DB.Where("id = ?", id).First(&webhook)
	if webhook.ID == 0 {
		Fail(c, "webhook不存在", gin.H{
			"webhookId": webhookId,
		})
		return
	}
	db.DB.Delete(&webhook)
	log.Println("删除webhook:", id)
	Success(c, "删除成功", nil)
}
