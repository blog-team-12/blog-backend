package database

// JwtBlacklist JWT 黑名单表
type JwtBlacklist struct {
	MODEL
	JWT string `json:"jwt" gorm:"type:text"` // Jwt
}
