/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/10 15:31
 * 描述     ：管理管理员账号
 */

package manage

import (
	"Evo/auth"
	"Evo/db"
	"Evo/model"
	"Evo/util"
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
		util.Fail(c, "参数绑定错误", gin.H{
			"admin": form,
		})
		return
	}

	var admin model.Admin
	db.DB.Where("name = ?", form.Name).First(&admin)
	if admin.ID != 0 {
		util.Fail(c, "账号已存在", gin.H{
			"admin": form,
		})
	}
	admin.Name = form.Name
	admin.Pwd = auth.Encode(form.Pwd)
	admin.Pwd = auth.Encode(admin.Pwd)
	db.DB.Create(&admin)
	admin.Pwd = ""
	log.Println("Add admin:", admin.Name)
	util.Success(c, "创建成功", nil)

}

type putAccountForm struct {
	AdminId uint   `binding:"required"`
	Name    string `binding:"required"`
	Pwd     string `binding:"required"`
}

// PutAccount 修改管理员账号
func PutAccount(c *gin.Context) {
	var form putAccountForm
	err := c.ShouldBind(&form)
	if err != nil {
		util.Fail(c, "参数绑定失败", gin.H{
			"admin": form,
		})
		return
	}
	var admin model.Admin
	db.DB.Select([]string{"id", "name", "pwd"}).First(&admin, form.AdminId)
	if admin.ID == 0 {
		util.Fail(c, "账号不存在", nil)
		return
	}
	if form.Name != "" {
		admin.Name = form.Name
	}
	if form.Pwd != "" {
		admin.Pwd = auth.Encode(form.Pwd)
	}

	db.DB.Model(&model.Admin{}).Where("id = ?", admin.ID).Updates(map[string]interface{}{
		"name": admin.Name,
		"pwd":  admin.Pwd,
	},
	)

	if err != nil {
		// 判断是不是查不到
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Fail(c, "账号不存在", gin.H{
				"admin": form,
			})
		} else {
			log.Println(err)
			util.Error(c, "修改失败", nil)
		}
		return
	}
	log.Println("modify account:", admin.Name)
	admin.Pwd = ""
	util.Success(c, "修改成功", gin.H{
		"admin": admin,
	})
}

// DelAccount 删除管理员账号
func DelAccount(c *gin.Context) {
	adminId := c.Query("adminId")
	id, err := strconv.Atoi(adminId)

	if err != nil {
		util.Fail(c, "参数错误", gin.H{
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
		util.Success(c, "success", nil)
		return
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Fail(c, "账号不存在", gin.H{
				"adminId": adminId,
			})
			return
		} else {
			log.Println(err)
			util.Error(c, "删除失败", gin.H{
				"adminId": adminId,
			})
			return
		}
	}
}

// GetAccount 获得所有管理员账号
func GetAccount(c *gin.Context) {
	admins := make([]model.Admin, 0)
	db.DB.Select([]string{"id", "created_at", "name"}).Find(&admins)

	util.Success(c, "success", gin.H{
		"admins": admins,
	})
}
