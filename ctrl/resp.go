/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间  :2022/07/09 14:18
 * 描述		:这里封装返沪请求的方法
 */

package ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const BindError ="参数绑定错误"

func Response(c *gin.Context, httpStatus int, code int, msg string, data gin.H) {
	c.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

func Fail(c *gin.Context, msg string, data gin.H) {
	Response(c, http.StatusOK, 400, msg, data)
}

func Success(c *gin.Context, msg string, data gin.H) {
	Response(c, http.StatusOK, 200, msg, data)
}

func Error(c *gin.Context, msg string, data gin.H) {
	Response(c, http.StatusOK, 500, msg, data)
}
