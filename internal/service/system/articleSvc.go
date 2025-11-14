package system

import (
	"context"
	"fmt"
	"personal_blog/global"
	"personal_blog/internal/model/consts"
	"personal_blog/internal/model/dto/request"
	esModel "personal_blog/internal/model/elasticsearch"
	"personal_blog/internal/repository"
	"personal_blog/internal/repository/interfaces"
	"personal_blog/pkg/articleUtils"
	esUtil "personal_blog/pkg/elasticSearch"
	"personal_blog/pkg/imageUtils"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ArticleSvc 文章服务
type ArticleSvc struct {
	articleRepo interfaces.ArticleRepository
}

// NewArticleSvc 创建文章服务实例
func NewArticleSvc(group *repository.Group) *ArticleSvc {
	return &ArticleSvc{
		articleRepo: group.SystemRepositorySupplier.GetArticleRepository(),
	}
}

func (a *ArticleSvc) ArticleCreate(ctx context.Context, req *request.ArticleCreateReq) error {
	// 1、通过关键字判断文章是存在
	b, err := articleUtils.Exists(ctx, req.Title)
	if err != nil {
		global.Log.Warn("检查标题失败，继续创建", zap.String("title", req.Title), zap.Error(err))
	}
	if b {
		global.Log.Warn("文章标题重复，继续创建", zap.String("title", req.Title))
	}
	// 2、设置结构体
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToCreate := &esModel.Article{
		CreatedAt:    now,
		UpdatedAt:    now,
		Cover:        req.Cover,
		Title:        req.Title,
		Keyword:      req.Title,
		Category:     req.Category,
		Tags:         req.Tags,
		Abstract:     req.Abstract,
		Content:      req.Content,
		VisibleRange: req.VisibleRange,
	}
	// 3、在事物中创建文章，并更新相关消息
	return a.articleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// 3.a 更新文章中的分类：新建文章仅需对新分类+1或创建
		if articleToCreate.Category != "" {
			if err = a.articleRepo.IncOrCreateCategory(ctx, tx, articleToCreate.Category); err != nil {
				global.Log.Error("更新分类计数失败", zap.String("category", articleToCreate.Category), zap.Error(err))
				return fmt.Errorf("更新分类计数失败: %v", err)
			}
		}
		// 3.b 更新标签计数
		if err = a.articleRepo.AddOrIncTag(ctx, tx, articleToCreate.Tags); err != nil {
			global.Log.Error("更新标签计数失败", zap.Strings("tags", articleToCreate.Tags), zap.Error(err))
			return fmt.Errorf("更新标签计数失败: %v", err)
		}
		// 3.c 创建文章
		err = articleUtils.CreateArticle(ctx, articleToCreate)
		if err != nil {
			global.Log.Error("创建文章失败", zap.String("title", articleToCreate.Title), zap.Error(err))
			return fmt.Errorf("创建文章失败: %v", err)
		}
		return nil
	})
}

// ArticleDelete 删除文章
func (a *ArticleSvc) ArticleDelete(ctx context.Context, req *request.ArticleDeleteReq) error {
	// 1、无数据，直接返回
	if len(req.IDs) == 0 {
		return nil
	}
	// 2、开启事物
	return a.articleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// 3、逐个删除
		for _, id := range req.IDs {
			// 3.a 获取文章
			articleToDelete, err := esUtil.Get(ctx, id)
			if err != nil {
				global.Log.Warn("获取文章失败",
					zap.String("id", id), zap.Error(err))
				return fmt.Errorf("获取文章失败: %v", err)
			}
			// 3.b 删除文章类别
			if err = a.articleRepo.DecOrDeleteCategory(ctx, tx, articleToDelete.Category); err != nil {
				global.Log.Error("更新分类计数失败",
					zap.String("category", articleToDelete.Category), zap.Error(err))
				return fmt.Errorf("更新分类计数失败: %v", err)
			}
			// 3.c 删除标签
			if err = a.articleRepo.DecOrDeleteTag(ctx, tx, articleToDelete.Tags); err != nil {
				global.Log.Error("更新标签计数失败",
					zap.Strings("tags", articleToDelete.Tags), zap.Error(err))
				return fmt.Errorf("更新标签计数失败: %v", err)
			}
			// 3.d 修改所有图片
			// 3.d.1 初始化图片封面类别
			imageSlice := []string{articleToDelete.Cover}
			if err = imageUtils.InitImagesCategory(ctx, tx, imageSlice); err != nil {
				global.Log.Error("初始化图片封面类别失败",
					zap.Strings("urls", imageSlice), zap.Error(err))
				return fmt.Errorf("初始化图片封面类别失败: %v", err)
			}
			// 3.d.2 获取所有插图
			imageSlice, err = imageUtils.FindIllustrations(articleToDelete.Content)
			if err != nil {
				global.Log.Warn("解析插图失败",
					zap.String("id", id), zap.Error(err))
				return fmt.Errorf("解析插图失败: %v", err)
			}
			// 3.d.3 修改所有插图类别
			if err = imageUtils.ChangeImagesCategory(ctx, tx, imageSlice, consts.Category(0)); err != nil {
				global.Log.Error("修改插图类别失败",
					zap.Strings("urls", imageSlice), zap.Error(err))
				return fmt.Errorf("修改插图类别失败: %v", err)
			}
			// 3.d
			// 同时删除所有评论 todo
			// 3.e 删除文章
			var ids = []string{id}
			if err = esUtil.Delete(ctx, ids); err != nil {
				global.Log.Error("删除文章失败",
					zap.String("id", id), zap.Error(err))
				return fmt.Errorf("删除文章失败: %v", err)
			}
		}
		return nil
	})
}
