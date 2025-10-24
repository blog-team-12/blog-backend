package consts

// 本包，均为后期拓展使用

// Storage 图片存储类型
type Storage int

const (
	Local Storage = iota // 本地
	Qiniu                // 七牛云
)
