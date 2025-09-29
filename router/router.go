package router

import (
	"github.com/gin-gonic/gin"
	"personal_blog/router/system"
)

type Routers struct {
	System system.RouterGroup
}

var GroupApp = new(Routers)

func InitRouter() *gin.Engine {
	Router := gin.New()

	// todo 后期填写logConfig
	Router.Use(gin.LoggerWithConfig(gin.LoggerConfig{}), gin.Recovery())
	systemRouter := GroupApp.System

	OutRouter := Router.Group("")
	systemRouter.InitRefreshTokenRouter(OutRouter)

	return Router
}
