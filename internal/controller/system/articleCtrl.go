package system

import (
	"personal_blog/global"
	"personal_blog/internal/model/dto/request"
	serviceSystem "personal_blog/internal/service/system"
	"personal_blog/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ArticleCtrl 文章控制器
type ArticleCtrl struct {
	articleSvc *serviceSystem.ArticleSvc
}

func (a *ArticleCtrl) CreateArticle(ctx *gin.Context) {
	var req request.ArticleCreateReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error("绑定数据错误", zap.Error(err))
		response.NewResponse[any, any](ctx).
			SetCode(global.StatusBadRequest).
			Failed("绑定数据错误", nil)
		return
	}
	err = a.articleSvc.ArticleCreate(ctx, &req)
	if err != nil {
		global.Log.Error("创建文章失败", zap.String("title", req.Title), zap.Error(err))
		response.NewResponse[any, any](ctx).
			SetCode(global.StatusInternalServerError).
			Failed("创建文章失败", nil)
		return
	}
	// 发布牛
	response.NewResponse[any, any](ctx).
		SetCode(global.StatusOK).
		Success("文章已创建", map[string]any{
			"title":         req.Title,
			"category":      req.Category,
			"tags":          req.Tags,
			"visible_range": req.VisibleRange,
		})
}
func (a *ArticleCtrl) DeleteArticle(ctx *gin.Context) {
	var req request.ArticleDeleteReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error("绑定数据错误", zap.Error(err))
		response.NewResponse[any, any](ctx).
			SetCode(global.StatusBadRequest).
			Failed("绑定数据错误", nil)
		return
	}
	err = a.articleSvc.ArticleDelete(ctx, &req)
	if err != nil {
		global.Log.Error("删除文章失败", zap.Strings("ids", req.IDs), zap.Error(err))
		response.NewResponse[any, any](ctx).
			SetCode(global.StatusInternalServerError).
			Failed("删除文章失败", nil)
		return
	}
	response.NewResponse[any, any](ctx).
		SetCode(global.StatusOK).
		Success("删除成功", map[string]any{
			"count": len(req.IDs),
		})
}
