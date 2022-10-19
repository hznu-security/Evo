/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/10 15:31
 * 描述     ：管理管理员账号
 */

package ctrl

import (
	"Evo/auth"
	"Evo/db"
	"Evo/model"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
)

// PostAccount 添加管理员账号
func PostAccount(c *gin.Context) {
	var form model.Admin
	err := c.ShouldBind(&form)
	if err != nil {
		Fail(c, "参数绑定错误", gin.H{
			"admin": form,
		})
		return
	}

	var admin model.Admin
	db.DB.Where("name = ?", form.Name).First(&admin)
	if admin.ID != 0 {
		Fail(c, "账号已存在", gin.H{
			"admin": form,
		})
	}
	admin.Name = form.Name
	admin.Pwd = auth.Encode(form.Pwd)
	admin.Pwd = auth.Encode(admin.Pwd)
	db.DB.Create(&admin)
	admin.Pwd = ""
	log.Println("Add admin:", admin.Name)
	Success(c, "创建成功", nil)

}

type putAccountForm struct {
	AdminId uint   `binding:"required"`
	Name    string `binding:"required"`
	Pwd     string `binding:"required"`
}

// PutAccount 修改管理员账号的密码
func PutAccount(c *gin.Context) {
	var form putAccountForm
	err := c.ShouldBind(&form)
	if err != nil {
		Fail(c, "参数绑定失败", gin.H{
			"admin": form,
		})
		return
	}

	var admin model.Admin
	// 事务写法
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := db.DB.First(&admin, form.AdminId).Error; err != nil {
			return err
		}
		admin.Pwd = auth.Encode(form.Pwd)
		if err = db.DB.Save(&admin).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// 判断是不是查不到
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(c, "账号不存在", gin.H{
				"admin": form,
			})
		} else {
			log.Println(err)
			Error(c, "修改失败", nil)
		}
		return
	}
	log.Println("add account:", admin.Name)
	Success(c, "修改成功", nil)
}

// DelAccount 删除管理员账号
func DelAccount(c *gin.Context) {
	adminId := c.Query("adminId")
	id, err := strconv.Atoi(adminId)

	if err != nil {
		Fail(c, "参数错误", gin.H{
			"adminId": adminId,
		})
		return
	}

	var admin model.Admin
	// 不需要写成事务
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&admin, id).Error; err != nil {
			return err
		}
		if err = tx.Delete(&admin).Error; err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		log.Println("delete account:", admin.Name)
		Success(c, "success", nil)
		return
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(c, "账号不存在", gin.H{
				"adminId": adminId,
			})
			return
		} else {
			log.Println(err)
			Error(c, "删除失败", gin.H{
				"adminId": adminId,
			})
			return
		}
	}

}

// GetAccount 获得所有管理员账号
func GetAccount(c *gin.Context) {
	admins := make([]model.Admin, 0)
	db.DB.Select([]string{"id", "name"}).Find(&admins)

	Success(c, "success", gin.H{
		"admins": admins,
	})
}
