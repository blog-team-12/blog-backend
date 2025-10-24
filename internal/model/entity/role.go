package entity

import "gorm.io/gorm"

// Role 角色表
type Role struct {
	gorm.Model
	Name   string `json:"name" gorm:"size:20;not null;comment:角色名称"`
	Code   string `json:"code" gorm:"size:20;not null;unique;comment:角色代码"`
	Status int    `json:"status" gorm:"default:1;comment:状态(1:启用 0:禁用)"`
	Desc   string `json:"desc" gorm:"size:200;comment:角色描述"`
	Menus  []Menu `json:"-" gorm:"many2many:role_menus;"`
}
