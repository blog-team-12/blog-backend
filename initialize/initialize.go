package initialize

import (
	"personal_blog/api"
	apiSystem "personal_blog/api/system"
	"personal_blog/core"
	"personal_blog/flag"
	"personal_blog/global"
	"personal_blog/service"
	serviceSystem "personal_blog/service/system"
)

func Init() {
	// 初始化配置
	core.InitConfig("config")
	// 初始化日志
	global.Log = core.InitLogger()
	// 为jwt黑名单开启本地存储
	core.OtherInit()
	// 连接数据库，初始化gorm
	global.DB = core.InitGorm()
	// 连接redis
	global.Redis = core.ConnectRedis()
	// 连接es
	global.ESClient = core.ConnectEs()

	// 开启flag
	flag.InitFlag()
	// 启动定时任务
	core.InitCron()
	// 加载jwt黑名单
	serviceSystem.LoadAll()
	
	// 业务函数--单例
	service.GroupApp = &service.Group{
		SystemServiceSupplier: serviceSystem.SetUp(),
	}
	// 控制函数
	api.ApiGroupApp = &api.ApiGroup{
		SystemApiGroup: apiSystem.SetUp(service.GroupApp),
	}

	// 开启函数
	core.RunServer()

}
