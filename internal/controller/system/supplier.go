package system

import (
	"personal_blog/internal/service"
)

type Supplier interface {
    GetRefreshTokenCtrl() *RefreshTokenCtrl
    GetBaseCtrl() *BaseCtrl
    GetUserCtrl() *UserCtrl
}

// SetUp 工厂函数-单例
func SetUp(service *service.Group) Supplier {
    cs := &controllerSupplier{}
    cs.refreshTokenCtrl = &RefreshTokenCtrl{
        jwtService: service.SystemServiceSupplier.GetJWTService(),
    }
    cs.baseCtrl = &BaseCtrl{
        baseService: service.SystemServiceSupplier.GetBaseService(),
    }
    cs.userCtrl = &UserCtrl{
        userService: service.SystemServiceSupplier.GetUserService(),
        jwtService:  service.SystemServiceSupplier.GetJWTService(),
    }
    return cs
}
