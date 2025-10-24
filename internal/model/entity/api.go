package entity

import "gorm.io/gorm"

// API 接口表
type API struct {
	gorm.Model
	Path    string `json:"path" gorm:"size:200;not null;comment:API路径"`
	Method  string `json:"method" gorm:"size:10;not null;comment:请求方法"`
	Detail  string `json:"detail" gorm:"size:100;comment:API描述"`
	GroupID uint   `json:"group_id" gorm:"index;comment:所属菜单组ID"`
	Status  int    `json:"status" gorm:"default:1;comment:状态(1:启用 0:禁用)"`
	Menu    Menu   `json:"-" gorm:"foreignKey:GroupID"`
}
