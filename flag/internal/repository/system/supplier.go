package system

import (
	"gorm.io/gorm"
	"personal_blog/internal/repository/adapter"
	"personal_blog/internal/repository/interfaces"
)

type Supplier interface {
	GetUserRepository() interfaces.UserRepository
	GetJWTRepository() interfaces.JWTRepository
}

// SetUp 工厂函数，统一管理 - 现在支持配置驱动
func SetUp(factoryConfig *adapter.FactoryConfig) Supplier {
	var userRepo interfaces.UserRepository
	var jwtRepo interfaces.JWTRepository

	switch factoryConfig.DatabaseType {
	case adapter.MySQL:
		if db, ok := factoryConfig.Connection.(*gorm.DB); ok {
			userRepo = NewUserRepository(db)
			jwtRepo = NewJwtRepository(db)
		}
	case adapter.MongoDB:
		// 未来可以添加Mongo	DB实现
		panic("MongoDB not implemented yet")
	case adapter.Redis:
		// 未来可以添加Redis实现
		panic("Redis not implemented yet")
	default:
		// 默认MySQL
		if db, ok := factoryConfig.Connection.(*gorm.DB); ok {
			userRepo = NewUserRepository(db)
			jwtRepo = NewJwtRepository(db)
		}
	}
	return &RepositorySupplier{
		userRepository: userRepo,
		jwtRepository:  jwtRepo,
	}
}
