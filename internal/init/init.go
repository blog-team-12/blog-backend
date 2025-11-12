package init

import (
    "context"
    "personal_blog/flag"
    "personal_blog/global"
    "personal_blog/internal/controller"
    apiSystem "personal_blog/internal/controller/system"
    "personal_blog/internal/core"
    "personal_blog/internal/repository"
    "personal_blog/internal/repository/adapter"
    "personal_blog/pkg/storage"
    _ "personal_blog/pkg/storage/local"
    _ "personal_blog/pkg/storage/qiniu"
    "time"

	"go.uber.org/zap"

	"personal_blog/internal/service"
	"personal_blog/internal/service/system"
)

func Init() {
    // 初始化配置
    core.InitConfig("configs")
    // 初始化存储驱动（根据全局配置选择 current 驱动）
    storage.InitFromConfig()
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
	// 初始化Casbin
	core.InitCasbin()

	// 开启flag
	flag.InitFlag()
	// 启动定时任务
	core.InitCron()

	// 初始化Repository层
	mysqlAdapter := &adapter.MySQLAdapter{}
	mysqlAdapter.SetConnection(global.DB) // 使用现有的数据库连接
	repository.InitRepositoryGroupWithAdapter(mysqlAdapter)

	// 加载jwt黑名单（使用Repository层）
	// 为初始化操作设置30秒超时，避免启动时卡死
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	system.LoadAllWithRepository(ctx, repository.GroupApp)

	// 业务函数--单例
	service.GroupApp = &service.Group{
		SystemServiceSupplier: system.SetUp(repository.GroupApp),
	}

	// 控制函数
	controller.ApiGroupApp = &controller.ApiGroup{
		SystemApiGroup: apiSystem.SetUp(service.GroupApp),
	}

	// 同步权限数据
	permissionService := service.GroupApp.SystemServiceSupplier.GetPermissionService()
	if err := permissionService.SyncAllPermissionsToCasbin(ctx); err != nil {
		global.Log.Error("权限同步失败", zap.Error(err))
	}

	// 开启函数
	core.RunServer()
}
