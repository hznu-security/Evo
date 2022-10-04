/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/10 10:51
 * 描述     ：管理员相关请求
 */

package ctrl

import (
	"Evo/auth"
	"Evo/db"
	"Evo/model"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminLogin 管理员登陆
func AdminLogin(c *gin.Context) {
	type loginForm struct {
		Name string `binding:"required" binding:"max=20"`
		Pwd  string `binding:"required" binding:"max=20"`
	}
	var form loginForm
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println(err.Error())
		Fail(c, "参数绑定失败", nil)
	}

	var admin model.Admin
	if err = db.DB.Where("name = ?", form.Name).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(c, "管理员不存在", nil)
			return
		} else {
			Error(c, "服务端错误", nil)
			log.Println(err.Error())
			return
		}
	}
	if !auth.Cmp(admin.Pwd, form.Pwd) {
		Fail(c, "密码错误", nil)
		return
	}
	token, err := auth.ReleaseToken(admin.ID, auth.ADMIN)
	if err != nil {
		log.Println(err.Error())
		Error(c, "服务端错误", nil)
		return
	}
	Success(c, "登陆成功", gin.H{
		"token": token,
	})
	log.Println(admin.Name, "login")
}
