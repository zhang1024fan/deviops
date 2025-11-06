package common

import (
	"dodevops-api/common/config"
	"log"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	redisConfig := config.GetRedisConfig()

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       0, // 使用默认数据库
	})

	// 测试连接
	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Printf("Redis connection failed: %v, will continue without cache", err)
		redisClient = nil
		return err
	}

	return nil
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return redisClient
}