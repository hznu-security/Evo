/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/09 14:14
 * 描述     ：程序入口，系统从这里启动
 */
package main

import (
	"Evo/config"
	"Evo/db"
	"Evo/router"
	"log"
)

func main() {
	config.InitConfig()
	db.InitDB()
	r := router.InitRouter()

	err := r.Run(":8080")
	if err != nil {
		log.Println(err.Error())
	}
}
