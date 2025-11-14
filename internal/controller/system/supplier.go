package system

import (
	"personal_blog/internal/service"
)

type Supplier interface {
	GetRefreshTokenCtrl() *RefreshTokenCtrl
	GetBaseCtrl() *BaseCtrl
	GetUserCtrl() *UserCtrl
	GetImageCtrl() *ImageCtrl
	GetArticleCtrl() *ArticleCtrl
}

// SetUp 工厂函数-单例
func SetUp(service *service.Group) Supplier {
	cs := &controllerSupplier{}
	cs.refreshTokenCtrl = &RefreshTokenCtrl{
		jwtService: service.SystemServiceSupplier.GetJWTSvc(),
	}
	cs.baseCtrl = &BaseCtrl{
		baseService: service.SystemServiceSupplier.GetBaseSvc(),
	}
	cs.userCtrl = &UserCtrl{
		userService: service.SystemServiceSupplier.GetUserSvc(),
		jwtService:  service.SystemServiceSupplier.GetJWTSvc(),
	}
	cs.imageCtrl = &ImageCtrl{
		imageService: service.SystemServiceSupplier.GetImageSvc(),
		jwtService:   service.SystemServiceSupplier.GetJWTSvc(),
	}
	cs.articleCtrl = &ArticleCtrl{
		articleSvc: service.SystemServiceSupplier.GetArticleSvc(),
	}
	return cs
}
