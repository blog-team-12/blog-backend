package router

import (
	"net/http"
	"os"
	"personal_blog/global"
	"personal_blog/internal/middleware"
	"personal_blog/internal/router/system"
	"personal_blog/internal/service"
	"strings"
	"time"

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
	// 应用核心中间件（日志、恢复、CORS）
	attachCoreMiddlewares(Router)
	// 将静态目录挂载到指定前缀
	mountStatic(Router)
	// 配置并挂载会话中间件
	attachSession(Router)
	// 超时中间间
	Router.Use(middleware.TimeoutMiddleware(30 * time.Second))

	systemRouter := GroupApp.System

	PublicGroup := Router.Group("")
	{
		// 刷新Token路由
		systemRouter.InitRefreshTokenRouter(PublicGroup)
		// 基础登录服务 - 获取验证码
		systemRouter.InitBaseRouter(PublicGroup)
		// 用户路由
		systemRouter.InitUserRouter(PublicGroup)
		// todo 登录、注册、健康检测.
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
		// 图片管理
		systemRouter.InitImageRouter(BusinessGroup)

		// 博客相关路由
	}
	return Router
}

// attachCoreMiddlewares 应用核心中间件（日志、恢复、CORS）
func attachCoreMiddlewares(r *gin.Engine) {
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	r.Use(middleware.CORSMiddleware())
}

// mountStatic 将静态目录挂载到指定前缀
func mountStatic(r *gin.Engine) {
	staticPrefix := strings.TrimSpace(global.Config.Static.Prefix)
	if staticPrefix == "" {
		staticPrefix = "/images"
	}
	staticPath := strings.TrimSpace(global.Config.Static.Path)
	if staticPath == "" {
		return
	}
	_ = os.MkdirAll(staticPath, 0755)
	r.StaticFS(staticPrefix, http.Dir(staticPath))
}

// attachSession 配置并挂载会话中间件
func attachSession(r *gin.Engine) {
	store := cookie.NewStore([]byte(global.Config.System.SessionsSecret))
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
	r.Use(sessions.Sessions("session", store))
}
