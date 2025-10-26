package system

import (
	"fmt"
	"go.uber.org/zap"
	"personal_blog/global"
	"personal_blog/internal/model/dto/request"
	serviceSystem "personal_blog/internal/service/system"
	"personal_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserCtrl struct {
	userService *serviceSystem.UserService
}

// Register 注册
func (u *UserCtrl) Register(ctx *gin.Context) {
	var req request.RegisterReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error("绑定数据错误",
			zap.Error(err))
		response.NewResponse[any, any](ctx).
			SetCode(global.StatusBadRequest).
			Failed(fmt.Sprintf("绑定数据错误: %v", err), nil)
		return
	}

	// 验证注册信息
	err = u.userService.VerifyRegister(ctx, &req)
	if err != nil {
		global.Log.Error("注册验证失败", zap.String("email", req.Email), zap.Error(err))
		response.NewResponse[any, any](ctx).SetCode(global.StatusBadRequest).
			Failed(fmt.Sprintf("注册验证失败: %v", err), nil)
		return
	}

	// 执行注册
	user, err := u.userService.Register(ctx, &req)
	if err != nil {
		global.Log.Error(
			"用户注册失败",
			zap.String("email", req.Email),
			zap.Error(err))
		response.NewResponse[any, any](ctx).
			SetCode(global.StatusInternalServerError).
			Failed(fmt.Sprintf("用户注册失败: %v", err), nil)
		return
	}

	global.Log.Info("用户注册成功",
		zap.String("email", req.Email),
		zap.Uint("userID", user.ID))
	response.NewResponse[any, any](ctx).SetCode(global.StatusOK).
		Success("注册成功", map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
		})

	// TODO: 注册成功后，生成 token 并返回
	//userApi.TokenNext(c, user)
}
