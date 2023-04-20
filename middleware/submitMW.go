/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:限制选手提交的中间件
 */

package middleware

import (
	"Evo/cron"
	"github.com/gin-gonic/gin"
)

func SubmitMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果比赛为开始或者比赛结束，阻止提交
		if cron.IsOn() {
			c.JSON(200, "比赛已结束,请勿提交")
			c.Abort()
			return
		}
	}
}
