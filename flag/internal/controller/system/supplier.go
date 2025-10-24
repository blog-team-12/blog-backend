package system

import (
	"personal_blog/internal/service"
)

type Supplier interface {
	GetRefreshTokenApi() *RefreshTokenApi
}

func SetUp(service *service.Group) Supplier {
	cs := &controllerSupplier{}
	cs.refreshTokenApi = &RefreshTokenApi{
		jwtService: service.SystemServiceSupplier.GetJWTService(),
	}
	return cs
}
