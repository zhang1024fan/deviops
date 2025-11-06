// 缓存redis配置初始化
// author xiaoRui

package redis

import (
	"context"
	"dodevops-api/common/config"
	"github.com/go-redis/redis/v8"
)

var (
	RedisDb *redis.Client //
)

// 初始化连接
func SetupRedisDb() error {
	var ctx = context.Background()
	RedisDb = redis.NewClient(&redis.Options{Addr: config.Config.Redis.Address, Password: config.Config.Redis.Password, DB: 0})
	_, err := RedisDb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
