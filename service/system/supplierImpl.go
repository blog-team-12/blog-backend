package system

// supplier implementation 用于底层实现

type serviceSupplier struct {
	esService  *EsService
	jwtService *JWTService
}

func (s *serviceSupplier) GetEsService() *EsService {
	return s.esService
}
func (s *serviceSupplier) GetJWTService() *JWTService {
	return s.jwtService
}
