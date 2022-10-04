/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/09 15:00
 * 描述     ：初始化数据库用的
 */

package db

import (
	"Evo/auth"
	"Evo/model"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	pwd := viper.GetString("datasource.pwd")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, pwd, host, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db

	// 初始化一个管理员进去
	var admin model.Admin
	db.Where("name = ?", "admin").First(&admin)
	if admin.ID == 0 {
		admin.Name = "admin"
		admin.Pwd = auth.Encode("123456")
		db.Create(&admin)
	}
	db.AutoMigrate(&model.Admin{}, &model.Attack{}, &model.Box{},
		&model.Challenge{}, &model.Flag{},
		&model.Notification{}, &model.Team{}, &model.Webhook{})
}
