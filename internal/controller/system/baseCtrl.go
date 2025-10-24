package system

import (
	"personal_blog/global"
	resp "personal_blog/internal/model/dto/response"
	serviceSystem "personal_blog/internal/service/system"
	"personal_blog/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type BaseCtrl struct {
	baseService *serviceSystem.BaseService
}

// 用来存储共享验证码
var store = base64Captcha.DefaultMemStore

// Captcha 生成数字验证码
func (b *BaseCtrl) Captcha(c *gin.Context) {
	helper := response.NewAPIHelper(c, "Captcha")

	// 调用服务层生成验证码，传递store
	id, b64s, err := b.baseService.GetCaptcha(store)
	if err != nil {
		helper.CommonError("Failed to generate captcha", global.StatusInternalServerError, err)
		response.NewResponse[resp.Captcha, resp.Captcha](c).
			SetCode(global.StatusInternalServerError).Failed("Failed to generate captcha", nil)
		return
	}

	// 成功响应
	response.NewResponse[resp.Captcha, resp.Captcha](c).
		SetCode(global.StatusOK).Success("验证码生成成功", resp.Captcha{
		CaptchaID: id,
		PicPath:   b64s,
	})
}
