package interfaces

import (
	"personal_blog/internal/model/entity"
	"time"
)

// JWTRepository JWT仓储接口
type JWTRepository interface {
	// Token黑名单管理
	AddToBlacklist(token string, expiry time.Time) error
	IsTokenBlacklisted(token string) (bool, error)
	CleanExpiredTokens() error

	// 用户Token记录
	SaveUserToken(userID uint, token string, expiry time.Time) error
	GetUserTokens(userID uint) ([]*entity.UserToken, error)
	RevokeUserToken(userID uint, token string) error
	RevokeAllUserTokens(userID uint) error

	// Token验证和刷新
	ValidateToken(token string) (bool, error)
	GetTokenInfo(token string) (*entity.TokenInfo, error)
	UpdateTokenExpiry(token string, newExpiry time.Time) error

	// 兼容现有Service层的方法
	CreateJwtBlacklist(jwtList *entity.JwtBlacklist) error
	IsJwtInBlacklist(jwt string) (bool, error)
	GetAllJwtBlacklist() ([]string, error)
	GetUserByID(id uint) (*entity.User, error)
}
