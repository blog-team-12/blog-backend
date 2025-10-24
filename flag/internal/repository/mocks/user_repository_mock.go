package mocks

import (
	"errors"
	"personal_blog/internal/model/entity"
	"personal_blog/internal/repository/interfaces"
)

// UserRepositoryMock 用户仓储Mock实现
type UserRepositoryMock struct {
	// 用于存储测试数据
	users map[uint]*entity.User
	// 用于控制方法行为
	shouldReturnError bool
	errorMessage      string
}

// NewUserRepositoryMock 创建Mock实例
func NewUserRepositoryMock() interfaces.UserRepository {
	return &UserRepositoryMock{
		users: make(map[uint]*entity.User),
	}
}

// SetError 设置是否返回错误（测试用）
func (m *UserRepositoryMock) SetError(shouldError bool, message string) {
	m.shouldReturnError = shouldError
	m.errorMessage = message
}

// AddTestUser 添加测试用户数据
func (m *UserRepositoryMock) AddTestUser(user *entity.User) {
	m.users[user.ID] = user
}

// GetByID Mock实现 - 从内存中获取
func (m *UserRepositoryMock) GetByID(id uint) (*entity.User, error) {
	if m.shouldReturnError {
		return nil, errors.New(m.errorMessage)
	}
	
	user, exists := m.users[id]
	if !exists {
		return nil, nil // 用户不存在
	}
	return user, nil
}

// GetByUsername Mock实现
func (m *UserRepositoryMock) GetByUsername(username string) (*entity.User, error) {
	if m.shouldReturnError {
		return nil, errors.New(m.errorMessage)
	}
	
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil
}

// GetByEmail Mock实现
func (m *UserRepositoryMock) GetByEmail(email string) (*entity.User, error) {
	if m.shouldReturnError {
		return nil, errors.New(m.errorMessage)
	}
	
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

// Create Mock实现 - 存储到内存
func (m *UserRepositoryMock) Create(user *entity.User) error {
	if m.shouldReturnError {
		return errors.New(m.errorMessage)
	}
	
	// 模拟自增ID
	if user.ID == 0 {
		user.ID = uint(len(m.users) + 1)
	}
	m.users[user.ID] = user
	return nil
}

// Update Mock实现
func (m *UserRepositoryMock) Update(user *entity.User) error {
	if m.shouldReturnError {
		return errors.New(m.errorMessage)
	}
	
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("用户不存在")
	}
	m.users[user.ID] = user
	return nil
}

// Delete Mock实现
func (m *UserRepositoryMock) Delete(id uint) error {
	if m.shouldReturnError {
		return errors.New(m.errorMessage)
	}
	
	delete(m.users, id)
	return nil
}

// GetUserList Mock实现
func (m *UserRepositoryMock) GetUserList(page, pageSize int) ([]*entity.User, int64, error) {
	if m.shouldReturnError {
		return nil, 0, errors.New(m.errorMessage)
	}
	
	var users []*entity.User
	for _, user := range m.users {
		users = append(users, user)
	}
	
	total := int64(len(users))
	
	// 简单分页逻辑
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(users) {
		return []*entity.User{}, total, nil
	}
	if end > len(users) {
		end = len(users)
	}
	
	return users[start:end], total, nil
}

// ExistsByUsername Mock实现
func (m *UserRepositoryMock) ExistsByUsername(username string) (bool, error) {
	if m.shouldReturnError {
		return false, errors.New(m.errorMessage)
	}
	
	for _, user := range m.users {
		if user.Username == username {
			return true, nil
		}
	}
	return false, nil
}

// ExistsByEmail Mock实现
func (m *UserRepositoryMock) ExistsByEmail(email string) (bool, error) {
	if m.shouldReturnError {
		return false, errors.New(m.errorMessage)
	}
	
	for _, user := range m.users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

// GetActiveUsers Mock实现
func (m *UserRepositoryMock) GetActiveUsers() ([]*entity.User, error) {
	if m.shouldReturnError {
		return nil, errors.New(m.errorMessage)
	}
	
	var activeUsers []*entity.User
	for _, user := range m.users {
		if !user.Freeze {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers, nil
}

// ValidateUser Mock实现
func (m *UserRepositoryMock) ValidateUser(username, password string) (*entity.User, error) {
	if m.shouldReturnError {
		return nil, errors.New(m.errorMessage)
	}
	
	user, err := m.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	
	// 简单密码验证（实际应该用bcrypt）
	if user.Password != password {
		return nil, errors.New("密码错误")
	}
	
	if user.Freeze {
		return nil, errors.New("用户已被冻结")
	}
	
	return user, nil
}

// UpdateLastLogin Mock实现
func (m *UserRepositoryMock) UpdateLastLogin(id uint) error {
	if m.shouldReturnError {
		return errors.New(m.errorMessage)
	}
	
	if _, exists := m.users[id]; !exists {
		return errors.New("用户不存在")
	}
	// Mock中不需要实际更新时间
	return nil
}