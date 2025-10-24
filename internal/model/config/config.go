package config

import "github.com/spf13/viper"

// Config 应用全局配置结构体，包含所有核心模块配置
type Config struct {
	ES      ES      `json:"es" yaml:"es"`           // Elasticsearch配置
	Redis   Redis   `json:"redis" yaml:"redis"`     // Redis配置
	Mysql   Mysql   `json:"mysql" yaml:"mysql"`     // MySQL数据库配置
	System  System  `json:"system" yaml:"system"`   // 系统服务配置
	Zap     Zap     `json:"zap" yaml:"zap"`         // 日志配置
	JWT     JWT     `json:"jwt" yaml:"jwt"`         // JWT认证配置
	Upload  Upload  `json:"upload" yaml:"upload"`   // 文件上传配置
	Captcha Captcha `json:"captcha" yaml:"captcha"` // 验证码配置
	Email   Email   `json:"email" yaml:"email"`     // 邮件发送配置
	Gaode   Gaode   `json:"gaode" yaml:"gaode"`     // 高德地图API配置
}

func NewConfig() *Config {
	// Elasticsearch配置初始化
	_es := &ES{
		URL:            viper.GetString("es.url"),
		Username:       viper.GetString("es.username"),
		Password:       viper.GetString("es.password"),
		IsConsolePrint: viper.GetBool("es.isConsolePrint"),
	}
	// Redis配置初始化
	_redis := &Redis{
		Address:  viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}
	// MySQL数据库配置初始化
	_mysql := &Mysql{
		Host:         viper.GetString("mysql.host"),
		Port:         viper.GetInt("mysql.port"),
		Config:       viper.GetString("mysql.config"),
		DBName:       viper.GetString("mysql.db_name"),
		Username:     viper.GetString("mysql.username"),
		Password:     viper.GetString("mysql.password"),
		MaxIdleConns: viper.GetInt("mysql.max_idle_conns"),
		MaxOpenConns: viper.GetInt("mysql.max_open_conns"),
		LogMode:      viper.GetString("mysql.log_mode"),
	}
	// 系统服务配置初始化
	_system := &System{
		Host:           viper.GetString("system.host"),
		Port:           viper.GetInt("system.port"),
		Env:            viper.GetString("system.env"),
		RouterPrefix:   viper.GetString("system.routerPrefix"),
		UseMultipoint:  viper.GetBool("system.useMultipoint"),
		SessionsSecret: viper.GetString("system.sessionsSecret"),
		OssType:        viper.GetString("system.ossType"),
	}
	// 日志配置初始化
	_zap := &Zap{
		Level:          viper.GetString("zap.level"),
		Filename:       viper.GetString("zap.filename"),
		MaxSize:        viper.GetInt("zap.max_size"),
		MaxBackups:     viper.GetInt("zap.max_backups"),
		MaxAge:         viper.GetInt("zap.max_age"),
		IsConsolePrint: viper.GetBool("zap.is_console_print"),
	}
	// JWT认证配置初始化
	_jwt := &JWT{
		AccessTokenSecret:      viper.GetString("jwt.access_token_secret"),
		RefreshTokenSecret:     viper.GetString("jwt.refresh_token_secret"),
		AccessTokenExpiryTime:  viper.GetString("jwt.access_token_expiry_time"),
		RefreshTokenExpiryTime: viper.GetString("jwt.refresh_token_expiry_time"),
		Issuer:                 viper.GetString("jwt.issuer"),
	}
	// 文件上传配置初始化
	_upload := &Upload{
		Size: viper.GetInt("upload.size"),
		Path: viper.GetString("upload.path"),
	}
	// 验证码配置初始化
	_captcha := &Captcha{
		Height:   viper.GetInt("captcha.height"),
		Width:    viper.GetInt("captcha.width"),
		Length:   viper.GetInt("captcha.length"),
		MaxSkew:  viper.GetFloat64("captcha.max_skew"),
		DotCount: viper.GetInt("captcha.dot_count"),
	}
	// 邮件发送配置初始化
	_email := &Email{
		Host:     viper.GetString("email.host"),
		Port:     viper.GetInt("email.port"),
		From:     viper.GetString("email.from"),
		Nickname: viper.GetString("email.nickname"),
		Secret:   viper.GetString("email.secret"),
		IsSSL:    viper.GetBool("email.is_ssl"),
	}
	// 高德地图API配置初始化
	_gaode := &Gaode{
		Enable: viper.GetBool("gaode.enable"),
		Key:    viper.GetString("gaode.key"),
	}

	return &Config{
		ES:      *_es,
		Redis:   *_redis,
		Mysql:   *_mysql,
		System:  *_system,
		Zap:     *_zap,
		JWT:     *_jwt,
		Upload:  *_upload,
		Captcha: *_captcha,
		Email:   *_email,
		Gaode:   *_gaode,
	}
}
