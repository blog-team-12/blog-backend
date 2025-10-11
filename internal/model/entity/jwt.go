package entity

import "time"

// UserToken 用户Token记录表
type UserToken struct {
	MODEL
	UserID    uint      `json:"user_id" gorm:"index"`                    // 用户ID
	User      User      `json:"user" gorm:"foreignKey:UserID"`           // 关联用户
	Token     string    `json:"token" gorm:"type:text;unique"`           // Token值
	TokenType string    `json:"token_type" gorm:"default:'access'"`      // Token类型：access, refresh
	ExpiresAt time.Time `json:"expires_at"`                              // 过期时间
	IsRevoked bool      `json:"is_revoked" gorm:"default:false"`         // 是否已撤销
	IP        string    `json:"ip"`                                      // 签发IP
	UserAgent string    `json:"user_agent"`                              // 用户代理
}

// TokenBlacklist Token黑名单表
type TokenBlacklist struct {
	MODEL
	Token     string    `json:"token" gorm:"type:text;unique"`  // Token值
	ExpiresAt time.Time `json:"expires_at"`                     // 过期时间
	Reason    string    `json:"reason"`                         // 加入黑名单原因
}

// TokenInfo Token信息结构体（用于返回）
type TokenInfo struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
	IsRevoked bool      `json:"is_revoked"`
}