package system

import (
	"github.com/gin-gonic/gin"
	"personal_blog/internal/controller"
)

type ArticleRouter struct {
}

func (ArticleRouter) InitArticleRouter(Router *gin.RouterGroup) {
	articleRouter := Router.Group("article")

	articleCtrl := controller.ApiGroupApp.SystemApiGroup.GetArticleCtrl()
	{
		articleRouter.POST("create", articleCtrl.CreateArticle)   // 创建文章
		articleRouter.DELETE("delete", articleCtrl.DeleteArticle) // 删除文章
		articleRouter.PUT("update", articleCtrl.ArticleUpdate)    // 更新文章
	}
}
