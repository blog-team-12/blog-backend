package system

// supplier implementation 用于底层实现
type serviceSupplier struct {
	esService         *EsService
	jwtService        *JWTService
	permissionService *PermissionService
	baseService       *BaseService
}

func (s *serviceSupplier) GetEsService() *EsService {
	return s.esService
}
func (s *serviceSupplier) GetJWTService() *JWTService {
	return s.jwtService
}
func (s *serviceSupplier) GetPermissionService() *PermissionService {
	return s.permissionService
}
func (s *serviceSupplier) GetBaseService() *BaseService {
	return s.baseService
}
