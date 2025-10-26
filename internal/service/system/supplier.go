package system

import "personal_blog/internal/repository"

type Supplier interface {
	GetEsService() *EsService
	GetJWTService() *JWTService
	GetPermissionService() *PermissionService
	GetBaseService() *BaseService
	GetUserService() *UserService
}

// SetUp 工厂函数，统一管理
func SetUp(repositoryGroup *repository.Group) Supplier {
	ss := &serviceSupplier{}
	ss.esService = NewEsService()
	ss.jwtService = NewJWTService(repositoryGroup)
	ss.permissionService = NewPermissionService(repositoryGroup)
	ss.baseService = NewBaseService() // 用不到repo层
	
	// UserService 需要依赖 PermissionService，所以在 permissionService 初始化后创建
	// 创建用户服务（注入权限服务）
	ss.userService = NewUserService(repositoryGroup, ss.permissionService)
	return ss
}
