package system

import "personal_blog/internal/repository"

type Supplier interface {
	GetEsService() *EsService
	GetJWTService() *JWTService
}

// SetUp 工厂函数，统一管理
func SetUp(repositoryGroup *repository.Group) Supplier {
	ss := &serviceSupplier{}
	ss.esService = NewEsService()
	ss.jwtService = NewJWTService(repositoryGroup)
	return ss
}
