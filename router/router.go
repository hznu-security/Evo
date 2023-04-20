package router

import (
	"Evo/ctrl/info"
	"Evo/ctrl/manage"
	"Evo/ctrl/team"
	"Evo/middleware"
	"Evo/starry"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	// 静态文件服务
	r.Static("/upload", "./upload")
	r.GET("/update", manage.UpdateScore)
	r.GET("/websocket", starry.ServeWebsocket)

	infoGroup := r.Group("/info")
	{
		infoGroup.GET("/time", info.Time)
		infoGroup.GET("/rank", info.GetRank)
	}
	r.GET("/time", info.Time)

	// 登录外所有接口都通过中间件进行验证
	manager := r.Group("/manager")
	{
		manager.POST("/checkdown", middleware.CheckAuth(), manage.Check)
		manager.POST("/login", manage.AdminLogin)
		manager.Use(middleware.AuthMW())
		manager.GET("/chart", manage.GetChart)
		account := manager.Group("/account")
		{
			account.POST("", manage.PostAccount)
			account.PUT("", manage.PutAccount)
			account.DELETE("", manage.DelAccount)
			account.GET("", manage.GetAccount)
		}
		flag := manager.Group("/flag")
		{
			//flag.POST("", manage.PostFlag)
			flag.GET("/generate", manage.GenerateFlag)
			flag.GET("/export", manage.ExportFlag)
			flag.POST("", manage.GetFlag)
		}
		config := manager.Group("/config")
		{
			config.GET("", manage.GetConfig)
			config.GET("/reset", manage.ResetConfig)
			config.PUT("", manage.PutConfig)
			config.GET("/start", manage.StartGame)
			config.GET("/terminate", manage.TerminateGame)
		}
		notification := manager.Group("/notification")
		{
			notification.PUT("", manage.PutNotice)
			notification.POST("", manage.PostNotice)
			notification.DELETE("", manage.DelNotice)
			notification.GET("", manage.GetNotice) //管理端获取通知
		}
		challenge := manager.Group("/challenge")
		{
			challenge.POST("", manage.PostChallenge)
			challenge.PUT("", manage.PutChallenge)
			challenge.POST("visible", manage.Visible)
			challenge.DELETE("", manage.DelChallenge)
			challenge.GET("", manage.GetChallenge)
		}
		box := manager.Group("/box")
		{
			box.POST("", manage.PostBox)
			box.PUT("", manage.PutBox)
			box.GET("", manage.GetBox)
			box.GET("/test", manage.TestSSH)
			box.DELETE("", manage.DelBox)
			box.GET("/reset", manage.ResetBox)
			box.GET("/testall", manage.TestSSHAll)
		}
		time := manager.Group("/team") //   8080:/manager/team/....
		{
			time.GET("", manage.GetTeam)
			time.POST("", manage.PostTeam)
			time.PUT("", manage.PutTeam)
			time.GET("/reset", manage.ResetPwd)
			time.DELETE("", manage.DelTeam)
			time.POST("/logo", manage.UploadLogo)
		}
		image := manager.Group("/image")
		{
			image.POST("", manage.PostImage)
			image.GET("", manage.GetImage)
			image.DELETE("", manage.DelImage)
		}
		starryGroup := manager.Group("/starry")
		{
			starryGroup.POST("/attack", starry.Attack)
			starryGroup.POST("/rank", starry.Rank)
			starryGroup.POST("/status", starry.Status)
			starryGroup.POST("/time", starry.Time)
			starryGroup.POST("/clear", starry.Clear)
			starryGroup.POST("/clearall", starry.ClearAll)
			starryGroup.POST("/round", starry.Round)
		}
	}

	teamGroup := r.Group("/team") //    8080:/team/.....   上下两个team不一样
	teamGroup.POST("/login", team.TeamLogin)
	{
		teamGroup.POST("/flag", middleware.SubmitMW(), middleware.AuthMW(), team.SubmitFlag) // 比赛结束后,不允许提交
		//team.Use(middleware.AuthMW())
		teamGroup.GET("/info", team.Info)
		teamGroup.GET("/gameBox", team.GetGameBox)
		teamGroup.GET("/notification", team.GetNotification) //选手端获取通知
	}
	return r
}
