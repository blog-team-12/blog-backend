package system

import "personal_blog/internal/repository/interfaces"

type RepositorySupplier struct {
	userRepository interfaces.UserRepository
	jwtRepository  interfaces.JWTRepository
	roleRepository interfaces.RoleRepository
	menuRepository interfaces.MenuRepository
	apiRepository  interfaces.APIRepository
}

func (r *RepositorySupplier) GetUserRepository() interfaces.UserRepository {
	return r.userRepository
}

func (r *RepositorySupplier) GetJWTRepository() interfaces.JWTRepository {
	return r.jwtRepository
}

func (r *RepositorySupplier) GetRoleRepository() interfaces.RoleRepository {
	return r.roleRepository
}

func (r *RepositorySupplier) GetMenuRepository() interfaces.MenuRepository {
	return r.menuRepository
}

func (r *RepositorySupplier) GetAPIRepository() interfaces.APIRepository {
	return r.apiRepository
}
