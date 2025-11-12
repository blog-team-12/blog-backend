package entity

import "personal_blog/internal/model/consts"

// Image 图片实体
type Image struct {
	MODEL
	Name     string          `json:"name" gorm:"type:varchar(128);not null;comment:'原始文件名或展示名称'"`
	URL      string          `json:"url" gorm:"type:varchar(512);uniqueIndex;not null;comment:'图片可访问URL（兼容更长域名/前缀）'"`
	Key      string          `json:"key" gorm:"type:varchar(256);index;comment:'存储Key（云存储/本地相对路径）'"`
	Category consts.Category `json:"category" gorm:"type:int;default:0;comment:'图片分类'"`
	Storage  consts.Storage  `json:"storage" gorm:"type:int;default:0;comment:'存储类型：本地/七牛'"`

	// 用户关联
	UserID *uint `json:"user_id" gorm:"index;comment:'所属用户ID'"`
	User   *User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
