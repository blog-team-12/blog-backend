package system

// supplier implementation 用于底层实现
type serviceSupplier struct {
	jwtService        *JWTService
	permissionService *PermissionService
	baseService       *BaseService
	userService       *UserService
	imageService      *ImageService
	articleSvc        *ArticleSvc
}

func (s *serviceSupplier) GetJWTSvc() *JWTService {
	return s.jwtService
}
func (s *serviceSupplier) GetPermissionSvc() *PermissionService {
	return s.permissionService
}
func (s *serviceSupplier) GetBaseSvc() *BaseService {
	return s.baseService
}
func (s *serviceSupplier) GetUserSvc() *UserService {
	return s.userService
}

func (s *serviceSupplier) GetImageSvc() *ImageService {
	return s.imageService
}

func (s *serviceSupplier) GetArticleSvc() *ArticleSvc {
	return s.articleSvc
}
