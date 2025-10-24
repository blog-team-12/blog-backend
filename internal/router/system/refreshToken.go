package system

import (
	"github.com/gin-gonic/gin"
	"personal_blog/internal/controller"
)

type refreshToken struct{}

func (r *refreshToken) InitRefreshTokenRouter(router *gin.RouterGroup) {
	refreshTokenRouter := router.Group("refreshToken")
	refreshTokenApi := controller.ApiGroupApp.SystemApiGroup.GetRefreshTokenCtrl()
	{
		// 刷新Api
		refreshTokenRouter.GET("", refreshTokenApi.RefreshToken)
	}

}
