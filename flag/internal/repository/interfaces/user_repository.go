package interfaces

import "personal_blog/internal/model/entity"

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基础CRUD操作
	GetByID(id uint) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uint) error

	// 业务相关查询
	GetUserList(page, pageSize int) ([]*entity.User, int64, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	GetActiveUsers() ([]*entity.User, error)

	// 认证相关
	ValidateUser(username, password string) (*entity.User, error)
	UpdateLastLogin(id uint) error
}
