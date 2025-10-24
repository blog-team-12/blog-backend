package consts

// Register 用户注册来源
type Register int

const (
	Email Register = iota // 邮箱验证码注册
	QQ                    // QQ登录注册
)
