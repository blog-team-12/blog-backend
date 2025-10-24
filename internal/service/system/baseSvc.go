package system

import (
	"github.com/mojocn/base64Captcha"
	"personal_blog/global"
)

type BaseService struct {
}

// NewBaseService 返回实例
func NewBaseService() *BaseService {
	return &BaseService{}
}

func (b *BaseService) GetCaptcha(store base64Captcha.Store) (string, string, error) {
	// 创建数字验证码的驱动
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)

	// 创建验证码对象
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验证码
	id, b64s, _, err := captcha.Generate()

	return id, b64s, err
}
