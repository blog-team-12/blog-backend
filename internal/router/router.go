package router

import (
	"personal_blog/global"
	"personal_blog/internal/middleware"
	"personal_blog/internal/router/system"
	"personal_blog/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

type Routers struct {
	System system.RouterGroup
}

var GroupApp = new(Routers)

func InitRouter() *gin.Engine {
	Router := gin.New()
	// 开启并记录堆栈信息
	Router.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	// 添加会话中间件
	var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret))
	Router.Use(sessions.Sessions("session", store))
	// 添加超时中间件
	Router.Use(middleware.TimeoutMiddleware(30)) // 30秒请求超时

	systemRouter := GroupApp.System
	// 公共路由 - 不需要认证
	PublicGroup := Router.Group("")
	{
		// 刷新Token路由
		systemRouter.InitRefreshTokenRouter(PublicGroup)
		systemRouter
		// todo 登录、注册、健康检测
	}

	// 系统管理路由 - 需要JWT认证与权限管理
	SystemGroup := Router.Group("")
	permissionMW := middleware.NewPermissionMiddleware(service.GroupApp) // 获取实例
	SystemGroup.Use(middleware.JWTAuth())                                // JWT认证
	SystemGroup.Use(permissionMW.CheckPermission())                      // 创建权限中间件
	{
		// 权限相关路由
	}
	// 业务路由组 - 需要JWT，但不需严格的权限控制
	BusinessGroup := Router.Group("")
	BusinessGroup.Use(middleware.JWTAuth())
	{
		// 博客相关路由
	}
	return Router
}
