package system

import (
	"personal_blog/internal/service"
)

type Supplier interface {
	GetRefreshTokenCtrl() *RefreshTokenCtrl
	GetBaseCtrl() *BaseCtrl
}

func SetUp(service *service.Group) Supplier {
	cs := &controllerSupplier{}
	cs.refreshTokenCtrl = &RefreshTokenCtrl{
		jwtService: service.SystemServiceSupplier.GetJWTService(),
	}
	cs.baseCtrl = &BaseCtrl{
		baseService: service.SystemServiceSupplier.GetBaseService(),
	}
	return cs
}
