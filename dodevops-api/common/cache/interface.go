package cache

import (
	"context"
	"time"
)

// ICacheService 缓存服务接口
type ICacheService interface {
	// 基础缓存操作
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (int64, error)
	
	// JSON 对象缓存操作
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, dest interface{}) error
	
	// 命名空间相关的缓存操作
	SetNamespaceList(ctx context.Context, clusterId uint, namespaces interface{}, expiration time.Duration) error
	GetNamespaceList(ctx context.Context, clusterId uint, dest interface{}) error
	SetNamespaceDetail(ctx context.Context, clusterId uint, namespaceName string, detail interface{}, expiration time.Duration) error
	GetNamespaceDetail(ctx context.Context, clusterId uint, namespaceName string, dest interface{}) error
	InvalidateNamespaceCache(ctx context.Context, clusterId uint, namespaceName ...string) error
}