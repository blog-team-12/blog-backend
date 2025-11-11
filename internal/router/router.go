package router

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"net/http"
	"personal_blog/global"
	"personal_blog/internal/middleware"
	"personal_blog/internal/router/system"
	"personal_blog/internal/service"
	"strings"
	"time"

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
	// 跨域中间件，开发环境允许 http://localhost:3000 访问
	Router.Use(middleware.CORSMiddleware())
	// 添加会话中间件
	var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret))
	// 根据环境切换会话 Cookie 选项：
	// - HTTP开发：SameSite=Lax, Secure=false（同源或顶层导航携带）
	// - HTTPS/生产：SameSite=None, Secure=true（允许跨站XHR/Fetch携带）
	var sameSite http.SameSite = http.SameSiteLaxMode
	var secure = false
	env := strings.ToLower(strings.TrimSpace(global.Config.System.Env))
	if env == "release" || strings.Contains(env, "https") {
		sameSite = http.SameSiteNoneMode
		secure = true
	}
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   5 * 60,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})
	Router.Use(sessions.Sessions("session", store))
	// 添加超时中间件
	Router.Use(middleware.TimeoutMiddleware(30 * time.Second)) // 30秒请求超时

	systemRouter := GroupApp.System

	PublicGroup := Router.Group("")
	{
		// 刷新Token路由
		systemRouter.InitRefreshTokenRouter(PublicGroup)
		// 基础登录服务 - 获取验证码
		systemRouter.InitBaseRouter(PublicGroup)
		// 用户路由
		systemRouter.InitUserRouter(PublicGroup)
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
