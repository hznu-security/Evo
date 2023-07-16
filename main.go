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
	"Evo/service/docker"
	"Evo/starry"
	"github.com/spf13/viper"
	"log"
)

func main() {
	// 初始化log
	//f, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	//if err != nil {
	//	return
	//}
	//defer func() {
	//	f.Close()
	//}()

	// 组合一下即可，os.Stdout代表标准输出流
	//multiWriter := io.MultiWriter(os.Stdout, f)
	//log.SetOutput(multiWriter)
	//
	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	config.InitConfig()
	db.InitDB()
	starry.Init()
	docker.InitDocker()
	r := router.InitRouter()

	port := viper.GetString("server.port")
	err := r.Run(":" + port)
	if err != nil {
		log.Println(err.Error())
	}
}
