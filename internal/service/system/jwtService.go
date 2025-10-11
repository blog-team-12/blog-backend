package system

import (
	"errors"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"personal_blog/global"
	"personal_blog/internal/model/dto/request"
	"personal_blog/internal/model/dto/response"
	"personal_blog/internal/model/entity"
	"personal_blog/internal/repository"
	erro "personal_blog/pkg/errors"
	"personal_blog/pkg/jwt"
	"personal_blog/pkg/util"
)

type JWTService struct {
	repositoryGroup *repository.Group
}

func NewJWTService(repositoryGroup *repository.Group) *JWTService {
	return &JWTService{
		repositoryGroup: repositoryGroup,
	}
}

// SetRedisJWT 将JWT存储到Redis中
func (j *JWTService) SetRedisJWT(jwt string, uuid uuid.UUID) error {
	// 解析配置中的JWT过期时间
	dr, err := util.ParseDuration(global.Config.JWT.RefreshTokenExpiryTime)
	if err != nil {
		return err
	}
	// 设置JWT在Redis中的过期时间
	return global.Redis.Set(uuid.String(), jwt, dr).Err()
}

// GetRedisJWT 从Redis中获取JWT
func (j *JWTService) GetRedisJWT(uuid uuid.UUID) (string, error) {
	// 从Redis获取指定uuid对应的JWT
	return global.Redis.Get(uuid.String()).Result()
}

// JoinInBlacklist 将JWT添加到黑名单
func (j *JWTService) JoinInBlacklist(jwtList entity.JwtBlacklist) error {
	// 将JWT记录插入到数据库中的黑名单表
	jwtRepo := j.repositoryGroup.SystemRepositorySupplier.GetJWTRepository()
	if err := jwtRepo.CreateJwtBlacklist(&jwtList); err != nil {
		return err
	}
	// 将JWT添加到内存中的黑名单缓存
	global.BlackCache.SetDefault(jwtList.JWT, struct{}{})
	return nil
}

// IsInBlacklist 检查JWT是否在黑名单中
func (j *JWTService) IsInBlacklist(jwt string) bool {
	// 从黑名单缓存中检查JWT是否存在
	_, ok := global.BlackCache.Get(jwt)
	return ok
}

// GetUserFromJWT 获取用户信息
func (j *JWTService) GetUserFromJWT(token string) (user *entity.User, jwtError *erro.JWTError) {
	jwtTool := jwt.NewJWT()
	refreshClaims, err := jwtTool.ParseAccessToken(token)
	if err != nil {
		return nil, erro.ClassifyJWTError(err)
		//return nil, fmt.Errorf("refresh Token is invalid %w", err)
	}
	// 验证用户是否存在，且未被冻结
	jwtRepo := j.repositoryGroup.SystemRepositorySupplier.GetJWTRepository()
	user, err = jwtRepo.GetUserByID(refreshClaims.UserID)
	if err != nil {
		return nil, erro.ClassifyJWTError(err)
	}
	if user.Freeze {
		return user, &erro.JWTError{
			Code:    global.StatusUserFrozen,
			Message: "用户已被冻结",
			Err:     errors.New("user has been frozen"),
		}
	}
	return user, nil
}

// GetAccessToken 获取
func (j *JWTService) GetAccessToken(token string) (resR *response.RefreshTokenResponse, jwtError *erro.JWTError) {
	user, jwtErr := j.GetUserFromJWT(token)
	if jwtErr != nil {
		return nil, jwtErr
	}
	jwtTool := jwt.NewJWT()
	claims := jwtTool.CreateAccessClaims(request.BaseClaims{
		UserID: user.ID,
		UUID:   user.UUID,
		RoleID: user.RoleID,
	})
	Token, err := jwtTool.CreateAccessToken(claims)
	if err != nil {
		return nil, &erro.JWTError{
			Code:    global.StatusInternalServerError,
			Message: "生成Token失败",
			Err:     errors.New("create token failed"),
		}
	}
	resR = &response.RefreshTokenResponse{
		AccessToken:          Token,
		AccessTokenExpiresAt: claims.ExpiresAt.Unix() * 1000,
	}
	return resR, nil
}

// LoadAll 从数据库加载所有的JWT黑名单并加入缓存
func LoadAll() {
	// 注意：这个函数在init阶段调用，此时Repository还未初始化
	// 所以暂时保持使用global.DB，后续可以考虑重构
	var data []string
	// 从数据库中获取所有的黑名单JWT
	if err := global.DB.Model(&entity.JwtBlacklist{}).Pluck("jwt", &data).Error; err != nil {
		// 如果获取失败，记录错误日志
		global.Log.Error("Failed to load JWT blacklist from the entity", zap.Error(err))
		return
	}
	// 将所有JWT添加到BlackCache缓存中
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	}
}

// LoadAllWithRepository 使用Repository加载所有的JWT黑名单并加入缓存
func LoadAllWithRepository(repositoryGroup *repository.Group) {
	jwtRepo := repositoryGroup.SystemRepositorySupplier.GetJWTRepository()
	data, err := jwtRepo.GetAllJwtBlacklist()
	if err != nil {
		// 如果获取失败，记录错误日志
		global.Log.Error("Failed to load JWT blacklist from the repository", zap.Error(err))
		return
	}
	// 将所有JWT添加到BlackCache缓存中
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	}
}
