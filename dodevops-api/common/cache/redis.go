package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// 缓存键前缀
	NamespaceListKeyPrefix   = "k8s:namespaces:list:"
	NamespaceDetailKeyPrefix = "k8s:namespaces:detail:"
	
	// 默认过期时间
	DefaultExpiration = 5 * time.Minute
	NamespaceListExpiration = 3 * time.Minute
	NamespaceDetailExpiration = 5 * time.Minute
)

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存服务
func NewRedisCache(client *redis.Client) ICacheService {
	return &RedisCache{
		client: client,
	}
}

// Set 设置缓存
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del 删除缓存
func (r *RedisCache) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查key是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (int64, error) {
	return r.client.Exists(ctx, key).Result()
}

// SetJSON 设置JSON对象缓存
func (r *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal json failed: %w", err)
	}
	
	return r.client.Set(ctx, key, jsonValue, expiration).Err()
}

// GetJSON 获取JSON对象缓存
func (r *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	jsonValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	
	return json.Unmarshal([]byte(jsonValue), dest)
}

// SetNamespaceList 设置命名空间列表缓存
func (r *RedisCache) SetNamespaceList(ctx context.Context, clusterId uint, namespaces interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("%s%d", NamespaceListKeyPrefix, clusterId)
	if expiration <= 0 {
		expiration = NamespaceListExpiration
	}
	return r.SetJSON(ctx, key, namespaces, expiration)
}

// GetNamespaceList 获取命名空间列表缓存
func (r *RedisCache) GetNamespaceList(ctx context.Context, clusterId uint, dest interface{}) error {
	key := fmt.Sprintf("%s%d", NamespaceListKeyPrefix, clusterId)
	return r.GetJSON(ctx, key, dest)
}

// SetNamespaceDetail 设置命名空间详情缓存
func (r *RedisCache) SetNamespaceDetail(ctx context.Context, clusterId uint, namespaceName string, detail interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("%s%d:%s", NamespaceDetailKeyPrefix, clusterId, namespaceName)
	if expiration <= 0 {
		expiration = NamespaceDetailExpiration
	}
	return r.SetJSON(ctx, key, detail, expiration)
}

// GetNamespaceDetail 获取命名空间详情缓存
func (r *RedisCache) GetNamespaceDetail(ctx context.Context, clusterId uint, namespaceName string, dest interface{}) error {
	key := fmt.Sprintf("%s%d:%s", NamespaceDetailKeyPrefix, clusterId, namespaceName)
	return r.GetJSON(ctx, key, dest)
}

// InvalidateNamespaceCache 使命名空间缓存失效
func (r *RedisCache) InvalidateNamespaceCache(ctx context.Context, clusterId uint, namespaceName ...string) error {
	var keysToDelete []string
	
	// 删除命名空间列表缓存
	listKey := fmt.Sprintf("%s%d", NamespaceListKeyPrefix, clusterId)
	keysToDelete = append(keysToDelete, listKey)
	
	// 删除指定命名空间的详情缓存
	if len(namespaceName) > 0 {
		for _, name := range namespaceName {
			detailKey := fmt.Sprintf("%s%d:%s", NamespaceDetailKeyPrefix, clusterId, name)
			keysToDelete = append(keysToDelete, detailKey)
		}
	} else {
		// 如果没有指定命名空间，删除该集群的所有命名空间详情缓存
		pattern := fmt.Sprintf("%s%d:*", NamespaceDetailKeyPrefix, clusterId)
		keys, err := r.client.Keys(ctx, pattern).Result()
		if err != nil {
			return fmt.Errorf("get keys by pattern failed: %w", err)
		}
		keysToDelete = append(keysToDelete, keys...)
	}
	
	if len(keysToDelete) > 0 {
		return r.Del(ctx, keysToDelete...)
	}
	
	return nil
}