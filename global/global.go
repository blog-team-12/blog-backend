package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"personal_blog/model/config"
)

var (
	Config     *config.Config             // 全局配置实例
	Log        *zap.Logger                // 全局日志实例
	DB         *gorm.DB                   // 全局数据库连接实例
	ESClient   *elasticsearch.TypedClient // 全局Elasticsearch客户端实例
	Redis      redis.Client               // 全局Redis客户端实例
	BlackCache local_cache.Cache          // 全局黑名单缓存实例
)
