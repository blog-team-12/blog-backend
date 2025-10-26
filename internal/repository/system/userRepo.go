package system

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"personal_blog/internal/model/entity"
	"personal_blog/internal/repository/interfaces"
)

// UserGormRepository 用户仓储GORM实现
type UserGormRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例，返回接口类型
func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserGormRepository{db: db}
}

// GetByID 根据ID获取用户
func (r *UserGormRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserGormRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserGormRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserGormRepository) Create(
	ctx context.Context,
	user *entity.User,
) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (r *UserGormRepository) Update(
	ctx context.Context,
	user *entity.User,
) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户（软删除）
func (r *UserGormRepository) Delete(
	ctx context.Context,
	id uint,
) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

// GetUserList 获取用户列表（分页）
func (r *UserGormRepository) GetUserList(
	ctx context.Context,
	page,
	pageSize int,
) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *UserGormRepository) ExistsByUsername(
	ctx context.Context,
	username string,
) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("username = ?", username).
		Count(&count).
		Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserGormRepository) ExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("email = ?", email).
		Count(&count).
		Error
	return count > 0, err
}

// GetActiveUsers 获取活跃用户（未冻结）
func (r *UserGormRepository) GetActiveUsers(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.WithContext(ctx).
		Where("freeze = ?", false).
		Find(&users).
		Error
	return users, err
}

// ValidateUser 验证用户登录
func (r *UserGormRepository) ValidateUser(
	ctx context.Context,
	username, password string,
) (*entity.User, error) {
	user, err := r.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	// 检查用户是否被冻结
	if user.Freeze {
		return nil, errors.New("用户已被冻结")
	}

	return user, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *UserGormRepository) UpdateLastLogin(
	ctx context.Context,
	id uint,
) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Update("updated_at", "NOW()").
		Error
}

// CheckEmailAddress 检查邮箱地址是否存在
func (r *UserGormRepository) CheckEmailAddress(
	ctx context.Context,
	email string,
) error {
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&entity.User{}).
		Error
	return err
}
