package system

import (
	"github.com/gin-gonic/gin"
	"personal_blog/api"
)

type refreshToken struct{}

func (r *refreshToken) InitRefreshTokenRouter(router *gin.RouterGroup) {
	refreshTokenRouter := router.Group("refreshToken")
	refreshTokenApi := api.ApiGroupApp.SystemApiGroup.GetRefreshTokenApi()
	{
		// 刷新Api
		refreshTokenRouter.GET("", refreshTokenApi.RefreshToken)
		//
	}

}
