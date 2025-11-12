package storage

import "personal_blog/global"

// 说明：
// - 为保持对外 API 稳定性，本文件保留“选择器”外观，但内部转发至全局管理器。
// - 全局管理器负责初始化所有驱动单例、线程安全的选择与切换。

// InitFromConfig 在应用启动阶段初始化所有驱动单例，并根据配置设置默认当前驱动。
// 配置项：storage.current 支持 local 或 qiniu；未配置时默认 local。
func InitFromConfig() {
	InitAll()
	if global.Config != nil {
		if name := global.Config.Storage.Current; name != "" {
			_ = SetCurrent(name)
		}
	}
}

// Current 获取当前驱动实例；若未初始化，返回 nil。
func Current() Driver { return CurrentDriver() }

// FromName 按名称获取指定驱动；未知名称时回退为当前驱动。
// 用途：允许在控制器中根据请求头/参数临时切换目标驱动（例如迁移压测）。
func FromName(name string) Driver {
	if d := DriverFromName(name); d != nil {
		return d
	}
	return CurrentDriver()
}
