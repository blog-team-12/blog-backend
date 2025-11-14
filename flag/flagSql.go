package flag

import (
    "personal_blog/global"
    "personal_blog/internal/model/entity"
)

// SQL 表结构迁移，如果表不存在，它会创建新表；如果表已经存在，它会根据结构更新表
func SQL() error {
	return global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&entity.User{},            // 用户表
		&entity.Login{},           // 登录日志表
		&entity.UserToken{},       // 用户Token记录表
		&entity.TokenBlacklist{},  // Token黑名单表
		&entity.JwtBlacklist{},    // JWT黑名单表（兼容现有代码）
		&entity.Role{},            // 角色表
		&entity.Menu{},            // 菜单表
		&entity.API{},             // API接口表
		&entity.UserRole{},        // 用户角色关联表
		&entity.API{},             // api表
		&entity.Image{},           // 图片表
		// Elasticsearch Article 非关系型存储，不参与 GORM 迁移
		&entity.ArticleCategory{}, // 文章类别表
		&entity.ArticleTag{},      // 文章标签表
		&entity.ArticleLike{},     // 文章点赞表
	)
}
