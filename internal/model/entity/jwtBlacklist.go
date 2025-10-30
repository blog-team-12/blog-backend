package entity

// JwtBlacklist JWT黑名单表 - 存储被撤销的JWT令牌，用于令牌失效管理
type JwtBlacklist struct {
	MODEL
	JWT string `json:"jwt" gorm:"type:text;not null;comment:'被撤销的JWT令牌'"`
	// 完整的JWT令牌字符串，用于验证时检查是否已被撤销
}
