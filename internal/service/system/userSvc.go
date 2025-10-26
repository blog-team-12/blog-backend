package system

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"personal_blog/global"
	"personal_blog/internal/model/consts"
	"personal_blog/internal/repository/interfaces"

	"personal_blog/internal/model/dto/request"
	"personal_blog/internal/model/entity"
	"personal_blog/internal/repository"
	"personal_blog/pkg/util"
	"time"
)

type UserService struct {
	userRepo          interfaces.UserRepository // 依赖接口而不是具体实现
	roleRepo          interfaces.RoleRepository // 角色仓储，用于获取默认角色
	permissionService *PermissionService        // 权限服务，用于RBAC角色分配
}

func NewUserService(repositoryGroup *repository.Group, permissionService *PermissionService) *UserService {
	return &UserService{
		userRepo:          repositoryGroup.SystemRepositorySupplier.GetUserRepository(),
		roleRepo:          repositoryGroup.SystemRepositorySupplier.GetRoleRepository(),
		permissionService: permissionService,
	}
}
func (u *UserService) VerifyRegister(ctx *gin.Context, req *request.RegisterReq) error {
	global.Log.Info("开始验证用户注册信息",
		zap.String("email", req.Email))

	session := sessions.Default(ctx)

	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil {
		global.Log.Error("会话中未找到邮箱信息",
			zap.String("requestEmail", req.Email))
		return errors.New("会话已过期，请重新发送验证码")
	}

	if savedEmail.(string) != req.Email {
		global.Log.Error("邮箱不一致",
			zap.String("sessionEmail", savedEmail.(string)),
			zap.String("requestEmail", req.Email))
		return errors.New("邮箱验证失败，请确认邮箱地址")
	}

	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil {
		global.Log.Error("会话中未找到验证码",
			zap.String("email", req.Email))
		return errors.New("验证码已过期，请重新发送")
	}

	if savedCode.(string) != req.VerificationCode {
		global.Log.Error("验证码不匹配",
			zap.String("email", req.Email),
			zap.String("inputCode", req.VerificationCode))
		return errors.New("验证码错误，请检查后重试")
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expire_time")
	if savedTime == nil {
		global.Log.Error("会话中未找到过期时间",
			zap.String("email", req.Email))
		return errors.New("验证码已过期，请重新发送")
	}

	if savedTime.(int64) < time.Now().Unix() {
		global.Log.Error("验证码已过期",
			zap.String("email", req.Email),
			zap.Int64("expireTime", savedTime.(int64)),
			zap.Int64("currentTime", time.Now().Unix()))
		return errors.New("验证码已过期，请重新发送")
	}

	global.Log.Info("用户注册信息验证成功", zap.String("email", req.Email))
	return nil
}

// Register 注册
func (u *UserService) Register(ctx *gin.Context, req *request.RegisterReq) (*entity.User, error) {
	// 检查邮箱是否已存在
	exists, err := u.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		global.Log.Error("检查邮箱是否存在时发生错误",
			zap.String("email", req.Email), zap.Error(err))
		return nil, errors.New("系统错误，请稍后重试")
	}
	if exists {
		global.Log.Error("邮箱已被注册",
			zap.String("email", req.Email))
		return nil, errors.New("该邮箱已被注册")
	}

	// 检查用户名是否已存在
	exists, err = u.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		global.Log.Error("检查用户名是否存在时发生错误",
			zap.String("username", req.Username),
			zap.Error(err))
		return nil, errors.New("系统错误，请稍后重试")
	}
	if exists {
		global.Log.Error("用户名已被注册",
			zap.String("username", req.Username))
		return nil, errors.New("该用户名已被注册")
	}

	// 创建用户实例（不直接设置RoleID，通过权限服务分配）
	user := &entity.User{
		Username: req.Username,
		Password: util.BcryptHash(req.Password),
		Email:    req.Email,
		UUID:     uuid.Must(uuid.NewV4()),
		Avatar:   "/image/avatar.jpg",
		Register: consts.Email,
		// 不直接设置 RoleID，将通过权限服务分配角色
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		global.Log.Error("创建用户失败",
			zap.String("email", req.Email),
			zap.String("username", req.Username),
			zap.Error(err))
		return nil, errors.New("创建用户失败，请稍后重试")
	}

	// 为新用户分配默认角色（从配置获取）
	defaultRoleCode := global.Config.System.DefaultRoleCode
	if defaultRoleCode == "" {
		defaultRoleCode = "user" // 兜底默认值
	}

	// 根据角色代码查找角色
	defaultRole, err := u.roleRepo.GetByCode(ctx, defaultRoleCode)
	if err != nil {
		global.Log.Error("获取默认角色失败",
			zap.String("role_code", defaultRoleCode),
			zap.Error(err))
		return nil, fmt.Errorf("获取默认角色失败: %w", err)
	}

	// 分配角色
	if err = u.permissionService.AssignRoleToUser(ctx, user.ID, defaultRole.ID); err != nil {
		global.Log.Error("分配默认角色失败",
			zap.Uint("user_id", user.ID),
			zap.Uint("role_id", defaultRole.ID),
			zap.String("role_code", defaultRole.Code),
			zap.Error(err))
		return nil, fmt.Errorf("分配默认角色失败: %w", err)
	}

	global.Log.Info("用户注册成功并分配默认角色",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("role_code", defaultRole.Code),
		zap.String("role_name", defaultRole.Name))
	return user, nil
}
