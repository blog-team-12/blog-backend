package entity

import (
	"github.com/gofrs/uuid"
	"personal_blog/internal/model/consts"
)

// User 用户表 - 存储系统用户的基本信息和状态
type User struct {
	MODEL
	UUID      uuid.UUID        `json:"uuid" gorm:"type:char(36);unique;not null;comment:'用户唯一标识符'"`                                    // 用户唯一标识符，用于对外暴露而不是数据库主键
	Username  string           `json:"username" gorm:"type:varchar(50);unique;not null;comment:'用户名'"`                                   // 用户登录名，必须唯一
	Password  string           `json:"-" gorm:"type:varchar(255);not null;comment:'用户密码哈希值'"`                                            // 用户密码的哈希值，不在JSON中返回
	Email     string           `json:"email" gorm:"type:varchar(100);unique;comment:'用户邮箱地址'"`                                           // 用户邮箱地址，用于登录和通知
	Openid    string           `json:"openid" gorm:"type:varchar(100);index;comment:'第三方登录OpenID'"`                                     // 第三方平台（如微信、QQ）的OpenID
	Avatar    string           `json:"avatar" gorm:"type:varchar(255);default:'';comment:'用户头像URL'"`                                     // 用户头像图片的URL地址
	Address   string           `json:"address" gorm:"type:varchar(200);default:'';comment:'用户地址信息'"`                                     // 用户的地理位置或地址信息
	Signature string           `json:"signature" gorm:"type:varchar(500);default:'签名是空白的，这位用户似乎比较低调。';comment:'用户个性签名'"`                  // 用户的个性签名或简介
	Register  consts.Register `json:"register" gorm:"type:tinyint;not null;default:1;comment:'注册来源'"`                                   // 用户注册来源（邮箱、第三方等）
	Freeze    bool             `json:"freeze" gorm:"type:boolean;not null;default:false;index;comment:'用户冻结状态'"`                         // 用户账户是否被冻结（禁用）
}
