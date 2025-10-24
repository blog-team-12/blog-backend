package response

// Captcha 验证码响应
type Captcha struct {
	CaptchaID string `json:"captcha_id"`
	PicPath   string `json:"pic_path"`
}

func (a Captcha) ToResponse(input *Captcha) *Captcha {
	return input
}
