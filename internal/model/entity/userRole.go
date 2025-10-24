package entity

import "gorm.io/gorm"

// UserRole 用户角色关联表
type UserRole struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"not null;comment:用户ID"`
	RoleID uint `json:"role_id" gorm:"not null;comment:角色ID"`
}
