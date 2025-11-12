package system

import (
    "github.com/gin-gonic/gin"
    "personal_blog/internal/controller"
)

type ImageRouter struct{}

func (i *ImageRouter) InitImageRouter(router *gin.RouterGroup) {
    imageRouter := router.Group("image")
    imageCtrl := controller.ApiGroupApp.SystemApiGroup.GetImageCtrl()
    {
        imageRouter.POST("upload", imageCtrl.Upload) // 上传图片
        imageRouter.GET("list", imageCtrl.List)      // 列出当前用户图片
        imageRouter.DELETE("/:id", imageCtrl.Delete) // 删除图片（修正路径语法）
    }
}