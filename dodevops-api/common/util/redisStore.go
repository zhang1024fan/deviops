package util

import (
	"context"
	"dodevops-api/common/constant"
	"dodevops-api/pkg/redis"
	"log"
	"time"
)

var ctx = context.Background()

type RedisStore struct{}

// Set 实现 base64Captcha.Store 接口方法，必须返回 error 类型
func (r RedisStore) Set(id string, value string) error {
	key := constant.LOGIN_CODE + id
	err := redis.RedisDb.Set(ctx, key, value, time.Minute*5).Err()
	if err != nil {
		log.Println("Redis Set Error:", err)
		return err
	}
	return nil
}

// Get 获取验证码
func (r RedisStore) Get(id string, clear bool) string {
	key := constant.LOGIN_CODE + id
	val, err := redis.RedisDb.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

// Verify 验证码校验
func (r RedisStore) Verify(id, answer string, clear bool) bool {
	v := r.Get(id, clear)
	return v == answer
}

// 获取redis中存放的k8s集群信息
func (r RedisStore) GetKubeCluster(name string) string {
	key := constant.KUBE_CLUSTER_CODE + name
	val, err := redis.RedisDb.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

// 插入redis中存放的k8s集群信息
func (r RedisStore) SetKubeCluster(name string, value string) bool {
	key := constant.KUBE_CLUSTER_CODE + name
	err := redis.RedisDb.Set(ctx, key, value, time.Hour).Err()
	if err != nil {
		log.Println("Redis Set Error:", err)
		return false
	}
	return true
}

// 插入redis中存放的k8s集群信息
func (r RedisStore) DelKubeCluster(name string) bool {
	key := constant.KUBE_CLUSTER_CODE + name
	err := redis.RedisDb.Del(ctx, key).Err()
	if err != nil {
		log.Println("Redis del Error:", err)
		return false
	}
	return true
}
func (r RedisStore) GetKubeClusterResource(name string) string {
	key := constant.KUBE_CLUSTER_CACHE_CODE + name
	value, err := redis.RedisDb.Get(ctx, key).Result()
	if err != nil {
		log.Println("GetKubeClusterResource Error:", err)
		return ""
	}
	return value
}
func (r RedisStore) SetKubeClusterResource(name string, value []byte) bool {
	key := constant.KUBE_CLUSTER_CACHE_CODE + name
	err := redis.RedisDb.Set(ctx, key, value, time.Minute*10).Err()
	if err != nil {
		log.Println("SetKubeClusterResource Error:", err)
		return false
	}
	return true
}
