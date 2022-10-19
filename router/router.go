package router

import (
	"Evo/ctrl"
	"Evo/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 登录外所有接口都通过中间件进行验证
	manager := r.Group("/manager")
	{
		manager.POST("/login", ctrl.AdminLogin)
		manager.Use(middleware.AuthMW())
		account := manager.Group("/account")
		{
			account.POST("", ctrl.PostAccount)
			account.PUT("", ctrl.PutAccount)
			account.DELETE("", ctrl.DelAccount)
			account.GET("", ctrl.GetAccount)
		}
		flag := manager.Group("/flag")
		{
			flag.POST("", ctrl.PostFlag)
			flag.GET("/generate", ctrl.GenerateFlag)
			flag.GET("/export", ctrl.ExportFlag)
			flag.GET("/filter", ctrl.FilterFlag)
		}
		config := manager.Group("/config")
		{
			config.GET("", ctrl.GetConfig)
			config.GET("/reset", ctrl.ResetConfig)
			config.PUT("", ctrl.PutConfig)
			config.GET("/start", ctrl.StartGame)
			config.GET("/terminate", ctrl.TerminateGame)
		}
		notification := manager.Group("/notification")
		{
			notification.PUT("", ctrl.PutNotice)
			notification.POST("", ctrl.PostNotice)
			notification.DELETE("", ctrl.DelNotice)
			notification.GET("", ctrl.GetNotice) //管理端获取通知
		}
		challenge := manager.Group("/challenge")
		{
			challenge.POST("", ctrl.PostChallenge)
			challenge.PUT("", ctrl.PutChallenge)
			challenge.POST("visible", ctrl.SetVisible)
			challenge.POST("unvisible", ctrl.SetUnVisible)
			challenge.DELETE("/challenge", ctrl.DelChallenge)
			challenge.GET("/challenge", ctrl.GetChallenge)
		}
		box := manager.Group("/box")
		{
			box.POST("", ctrl.PostBox)
			box.PUT("", ctrl.PutBox)
			box.GET("", ctrl.GetBox)
			box.GET("", ctrl.GenerateFlag)
			box.GET("", ctrl.TestSSH)
			box.DELETE("", ctrl.DelBox)
			box.GET("/reset", ctrl.ResetBox)
		}
		team := manager.Group("/team") //   8080:/manager/team/....
		{
			team.GET("", ctrl.GetTeam)
			team.POST("", ctrl.PostTeam)
			team.PUT("", ctrl.PutTeam)
			team.GET("/reset", ctrl.ResetPwd)
			team.DELETE("", ctrl.DelTeam)
			team.POST("/logo", ctrl.UploadLogo)
		}
		image := manager.Group("/image")
		image.POST("", ctrl.PostImage)
		image.GET("", ctrl.GetImage)
		image.DELETE("", ctrl.DelImage)
	}

	team := r.Group("/team") //    8080:/team/.....   上下两个team不一样
	team.POST("/login", ctrl.TeamLogin)
	{
		team.POST("/flag", middleware.SubmitMW(), middleware.AuthMW(), ctrl.SubmitFlag) // 比赛结束后,不允许提交
		team.Use(middleware.AuthMW())
		team.GET("/info", ctrl.Info)
		team.GET("/rank", ctrl.GetRank)
		team.GET("/notification", ctrl.GetNotification) //选手端获取通知
	}
	return r
}
