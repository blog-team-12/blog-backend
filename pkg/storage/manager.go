package storage

import (
	"sync"
)

// Manager 负责管理所有存储驱动的单例，并提供统一的选择与切换能力。
// 目标：
// - 在应用启动时初始化“所有驱动单例”（而非懒加载）；
// - 提供按名称获取指定驱动、获取当前驱动、设置当前驱动等能力；
// - 线程安全，供并发请求下的读取与选择；
// - 保持最小实现，后续可根据配置扩展（如从配置读取当前驱动、驱动参数等）。
type Manager struct {
	mu          sync.RWMutex
	drivers     map[string]Driver
	current     string
	initialized bool
}

var manager = &Manager{}

// RegisterDriver 由具体驱动在其包的 init() 中调用，注册单例实例。
// 注意：name 不得重复；若重复注册，将覆盖旧实例。
func RegisterDriver(name string, drv Driver) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if manager.drivers == nil {
		manager.drivers = make(map[string]Driver)
	}
	manager.drivers[name] = drv
}

// InitAll 标记管理器已初始化，并设置默认当前驱动（若尚未设置）。
// 具体驱动的实例由各自包在 init() 时注册，避免包循环依赖。
func InitAll() {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if manager.initialized {
		return
	}
	// 默认当前驱动为 local（后续可按配置注入）；若没有注册 local，则保持空，等待配置设置
	if manager.current == "" {
		manager.current = "local"
	}
	manager.initialized = true
}

// CurrentDriver 返回当前驱动实例；若尚未初始化，返回 nil。
func CurrentDriver() Driver {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	if !manager.initialized {
		return nil
	}
	return manager.drivers[manager.current]
}

// DriverFromName 按名称返回指定驱动；未知名称返回 nil。
// 用途：可在控制器中根据请求头/参数选择不同驱动进行操作。
func DriverFromName(name string) Driver {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	if !manager.initialized {
		return nil
	}
	return manager.drivers[name]
}

// SetCurrent 将当前驱动切换为指定名称；返回是否切换成功。
func SetCurrent(name string) bool {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if !manager.initialized {
		return false
	}
	if _, ok := manager.drivers[name]; !ok {
		return false
	}
	manager.current = name
	return true
}

// DriverNames 返回已注册的驱动名称列表。
func DriverNames() []string {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	names := make([]string, 0, len(manager.drivers))
	for k := range manager.drivers {
		names = append(names, k)
	}
	return names
}
