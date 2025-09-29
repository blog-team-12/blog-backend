package system

type Supplier interface {
	GetEsService() *EsService
	GetJWTService() *JWTService
}

// SetUp 工厂函数，统一管理
func SetUp() Supplier {
	ss := &serviceSupplier{}
	ss.esService = NewEsService()
	ss.jwtService = NewJWTService()
	return ss
}
