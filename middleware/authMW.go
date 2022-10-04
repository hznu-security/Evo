/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/10 17:44
 * 描述     ：鉴权中间件
 */

package middleware

import (
	"Evo/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}
		tokenStr = tokenStr[7:]
		//没带token
		token, claims, err := auth.ParseToken(tokenStr)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			// Abort 函数会终端这次请求，后面的函数不会再被调用
			c.Abort()
			return
		}
		//拿到url中的路径，判断是对选手端还是对管理端的请求
		url := c.Request.URL.Path
		var role uint8
		if strings.HasPrefix(url, "/team") {
			role = auth.TEAM
		} else {
			role = auth.ADMIN
		}
		if claims.Role != role {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
		}
		c.Set("teamId", claims.ID)
		c.Next()
	}
}
