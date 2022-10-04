package router

import (
	"Evo/ctrl"
	"Evo/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/manager/login", ctrl.AdminLogin)
	r.POST("/manager/team", middleware.AuthMW(), ctrl.PostTeam)
	r.POST("/manager/image", ctrl.PostImage)
	r.GET("/manager/image", ctrl.GetImage)
	r.DELETE("/manager/image", ctrl.DelImage)
	r.POST("/team/login", ctrl.TeamLogin)

	return r
}
