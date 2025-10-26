package core

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"personal_blog/global"
	"personal_blog/internal/model/config"
	"strings"
)

// InitConfig 初始化配置 - 以环境变量优先
func InitConfig(path string) {
	// 条件初始化
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 读取并监听配置
	viper.WatchConfig() // 监视配置文件的更改
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")
	viper.SetConfigName("configs")

	if err := viper.ReadInConfig(); err != nil {
		global.Log.Fatal("配置文件加载失败") // 早期处理错误
	}

	// 绑定验证码相关配置到环境变量
	_ = viper.BindEnv("captcha.height", "CAPTCHA_HEIGHT")
	_ = viper.BindEnv("captcha.width", "CAPTCHA_WIDTH")
	_ = viper.BindEnv("captcha.length", "CAPTCHA_LENGTH")
	_ = viper.BindEnv("captcha.max_skew", "CAPTCHA_MAX_SKEW")
	_ = viper.BindEnv("captcha.dot_count", "CAPTCHA_DOT_COUNT")

	// 绑定邮件相关配置到环境变量
	_ = viper.BindEnv("email.host", "EMAIL_HOST")
	_ = viper.BindEnv("email.port", "EMAIL_PORT")
	_ = viper.BindEnv("email.from", "EMAIL_FROM")
	_ = viper.BindEnv("email.nickname", "EMAIL_NICKNAME")
	_ = viper.BindEnv("email.secret", "EMAIL_SECRET")
	_ = viper.BindEnv("email.is_ssl", "EMAIL_IS_SSL")

	// 绑定Elasticsearch相关配置到环境变量（补充完整）
	_ = viper.BindEnv("es.url", "ES_URL")
	_ = viper.BindEnv("es.username", "ES_USERNAME")
	_ = viper.BindEnv("es.password", "ES_PASSWORD")
	_ = viper.BindEnv("es.is_console_print", "ES_IS_CONSOLE_PRINT")

	// 绑定高德地图相关配置到环境变量
	_ = viper.BindEnv("gaode.enable", "GAODE_ENABLE")
	_ = viper.BindEnv("gaode.key", "GAODE_KEY")

	// 绑定JWT相关配置到环境变量
	_ = viper.BindEnv("jwt.access_token_secret", "JWT_ACCESS_TOKEN_SECRET")
	_ = viper.BindEnv("jwt.refresh_token_secret", "JWT_REFRESH_TOKEN_SECRET")
	_ = viper.BindEnv("jwt.access_token_expiry_time", "JWT_ACCESS_TOKEN_EXPIRY_TIME")
	_ = viper.BindEnv("jwt.refresh_token_expiry_time", "JWT_REFRESH_TOKEN_EXPIRY_TIME")
	_ = viper.BindEnv("jwt.issuer", "JWT_ISSUER")

	// 绑定MySQL相关配置到环境变量（补充完整）
	_ = viper.BindEnv("mysql.host", "DB_HOST")
	_ = viper.BindEnv("mysql.port", "DB_PORT")
	_ = viper.BindEnv("mysql.configs", "DB_CONFIG")
	_ = viper.BindEnv("mysql.db_name", "DB_NAME")
	_ = viper.BindEnv("mysql.username", "DB_USERNAME")
	_ = viper.BindEnv("mysql.password", "DB_PASSWORD")
	_ = viper.BindEnv("mysql.max_idle_conns", "DB_MAX_IDLE_CONNS")
	_ = viper.BindEnv("mysql.max_open_conns", "DB_MAX_OPEN_CONNS")
	_ = viper.BindEnv("mysql.log_mode", "DB_LOG_MODE")

	// 绑定七牛云相关配置到环境变量
	_ = viper.BindEnv("qiniu.zone", "QINIU_ZONE")
	_ = viper.BindEnv("qiniu.bucket", "QINIU_BUCKET")
	_ = viper.BindEnv("qiniu.img_path", "QINIU_IMG_PATH")
	_ = viper.BindEnv("qiniu.access_key", "QINIU_ACCESS_KEY")
	_ = viper.BindEnv("qiniu.secret_key", "QINIU_SECRET_KEY")
	_ = viper.BindEnv("qiniu.use_https", "QINIU_USE_HTTPS")
	_ = viper.BindEnv("qiniu.use_cdn_domains", "QINIU_USE_CDN_DOMAINS")

	// 绑定QQ登录相关配置到环境变量
	_ = viper.BindEnv("qq.enable", "QQ_ENABLE")
	_ = viper.BindEnv("qq.app_id", "QQ_APP_ID")
	_ = viper.BindEnv("qq.app_key", "QQ_APP_KEY")
	_ = viper.BindEnv("qq.redirect_uri", "QQ_REDIRECT_URI")

	// 绑定Redis相关配置到环境变量（补充完整）
	_ = viper.BindEnv("redis.address", "REDIS_ADDRESS")
	_ = viper.BindEnv("redis.password", "REDIS_PASSWORD")
	_ = viper.BindEnv("redis.db", "REDIS_DB")

	// 绑定系统服务相关配置到环境变量（补充完整）
	_ = viper.BindEnv("system.host", "SYSTEM_HOST")
	_ = viper.BindEnv("system.port", "PORT") // 复用PORT环境变量
	_ = viper.BindEnv("system.env", "SYSTEM_ENV")
	_ = viper.BindEnv("system.router_prefix", "SYSTEM_ROUTER_PREFIX")
	_ = viper.BindEnv("system.use_multipoint", "SYSTEM_USE_MULTIPOINT")
	_ = viper.BindEnv("system.sessions_secret", "SYSTEM_SESSIONS_SECRET")
	_ = viper.BindEnv("system.oss_type", "SYSTEM_OSS_TYPE")

	// 绑定文件上传相关配置到环境变量
	_ = viper.BindEnv("upload.size", "UPLOAD_SIZE")
	_ = viper.BindEnv("upload.path", "UPLOAD_PATH")

	// 绑定网站信息相关配置到环境变量
	_ = viper.BindEnv("website.logo", "WEBSITE_LOGO")
	_ = viper.BindEnv("website.full_logo", "WEBSITE_FULL_LOGO")
	_ = viper.BindEnv("website.title", "WEBSITE_TITLE")
	_ = viper.BindEnv("website.slogan", "WEBSITE_SLOGAN")
	_ = viper.BindEnv("website.slogan_en", "WEBSITE_SLOGAN_EN")
	_ = viper.BindEnv("website.description", "WEBSITE_DESCRIPTION")
	_ = viper.BindEnv("website.version", "WEBSITE_VERSION")
	_ = viper.BindEnv("website.created_at", "WEBSITE_CREATED_AT")
	_ = viper.BindEnv("website.icp_filing", "WEBSITE_ICP_FILING")
	_ = viper.BindEnv("website.public_security_filing", "WEBSITE_PUBLIC_SECURITY_FILING")
	_ = viper.BindEnv("website.bilibili_url", "WEBSITE_BILIBILI_URL")
	_ = viper.BindEnv("website.gitee_url", "WEBSITE_GITEE_URL")
	_ = viper.BindEnv("website.github_url", "WEBSITE_GITHUB_URL")
	_ = viper.BindEnv("website.blog_url", "WEBSITE_BLOG_URL")
	_ = viper.BindEnv("website.name", "WEBSITE_NAME")
	_ = viper.BindEnv("website.job", "WEBSITE_JOB")
	_ = viper.BindEnv("website.address", "WEBSITE_ADDRESS")
	_ = viper.BindEnv("website.email", "WEBSITE_EMAIL")
	_ = viper.BindEnv("website.qq_image", "WEBSITE_QQ_IMAGE")
	_ = viper.BindEnv("website.wechat_image", "WEBSITE_WECHAT_IMAGE")

	// 绑定Zap日志相关配置到环境变量
	_ = viper.BindEnv("zap.level", "ZAP_LEVEL")
	_ = viper.BindEnv("zap.filename", "ZAP_FILENAME")
	_ = viper.BindEnv("zap.max_size", "ZAP_MAX_SIZE")
	_ = viper.BindEnv("zap.max_backups", "ZAP_MAX_BACKUPS")
	_ = viper.BindEnv("zap.max_age", "ZAP_MAX_AGE")
	_ = viper.BindEnv("zap.is_console_print", "ZAP_IS_CONSOLE_PRINT")

	global.Log.Info("--------- configs list--------\n")
	for _, key := range viper.AllKeys() {
		global.Log.Info("configs",
			zap.String("key", key),
			zap.Any("value", viper.Get(key)))
	}
	global.Log.Info("-----------------------------\n")
	// 传递到全局
	global.Config = config.NewConfig()
}
