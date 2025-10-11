package system

import "personal_blog/internal/repository/interfaces"

type RepositorySupplier struct {
	userRepository interfaces.UserRepository
	jwtRepository  interfaces.JWTRepository
}

func (r *RepositorySupplier) GetUserRepository() interfaces.UserRepository {
	return r.userRepository
}

func (r *RepositorySupplier) GetJWTRepository() interfaces.JWTRepository {
	return r.jwtRepository
}
