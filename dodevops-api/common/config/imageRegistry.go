// 镜像仓库配置管理
// author xiaoRui
package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// ImageRegistryConfig 镜像仓库配置
type ImageRegistryConfig struct {
	Registry struct {
		Private string `yaml:"private"` // 私有仓库地址
		Public  string `yaml:"public"`  // 公共仓库地址
	} `yaml:"registry"`
	Images map[string]map[string]string `yaml:"images"` // 镜像配置
}

var (
	imageConfig     *ImageRegistryConfig
	imageConfigOnce sync.Once
)

// GetImageRegistry 获取镜像仓库配置（单例）
func GetImageRegistry() *ImageRegistryConfig {
	imageConfigOnce.Do(func() {
		imageConfig = loadImageConfig()
	})
	return imageConfig
}

// loadImageConfig 加载镜像配置文件
func loadImageConfig() *ImageRegistryConfig {
	configFile := "common/config/images.yaml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		// 如果配置文件不存在，返回默认配置
		return getDefaultImageConfig()
	}

	var config ImageRegistryConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return getDefaultImageConfig()
	}

	return &config
}

// getDefaultImageConfig 获取默认镜像配置
func getDefaultImageConfig() *ImageRegistryConfig {
	return &ImageRegistryConfig{
		Registry: struct {
			Private string `yaml:"private"`
			Public  string `yaml:"public"`
		}{
			Private: "crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s",
			Public:  "docker.io",
		},
		Images: make(map[string]map[string]string),
	}
}

// GetImage 获取指定服务的镜像地址
// service: 服务名称（如: mysql, redis, nodejs）
// version: 版本号（如: 8.0, 6.2, 18）
// usePrivate: 是否使用私有仓库（默认true）
func (c *ImageRegistryConfig) GetImage(service, version string, usePrivate bool) string {
	// 从配置中查找镜像
	if versions, ok := c.Images[service]; ok {
		if imageTpl, ok := versions[version]; ok {
			// 替换模板中的 {{.Registry}}
			registry := c.Registry.Public
			if usePrivate {
				registry = c.Registry.Private
			}
			return replaceRegistry(imageTpl, registry)
		}
	}

	// 如果配置中没有，返回默认格式
	registry := c.Registry.Public
	if usePrivate {
		registry = c.Registry.Private
	}
	return fmt.Sprintf("%s/%s:%s", registry, service, version)
}

// GetPrivateRegistry 获取私有仓库地址
func (c *ImageRegistryConfig) GetPrivateRegistry() string {
	return c.Registry.Private
}

// GetPublicRegistry 获取公共仓库地址
func (c *ImageRegistryConfig) GetPublicRegistry() string {
	return c.Registry.Public
}

// replaceRegistry 替换镜像模板中的仓库地址
func replaceRegistry(template, registry string) string {
	// 替换 {{.Registry}} 占位符
	return strings.Replace(template, "{{.Registry}}", registry, -1)
}

// 便捷方法：获取常用镜像

// GetMySQLImage 获取MySQL镜像
func GetMySQLImage(version string) string {
	return GetImageRegistry().GetImage("mysql", version, true)
}

// GetRedisImage 获取Redis镜像
func GetRedisImage(version string) string {
	return GetImageRegistry().GetImage("redis", version, true)
}

// GetPostgreSQLImage 获取PostgreSQL镜像
func GetPostgreSQLImage(version string) string {
	return GetImageRegistry().GetImage("postgresql", version, true)
}

// GetNodeJSImage 获取NodeJS镜像
func GetNodeJSImage(version string) string {
	return GetImageRegistry().GetImage("nodejs", version, true)
}

// GetJavaImage 获取Java镜像
func GetJavaImage(version string) string {
	return GetImageRegistry().GetImage("java", version, true)
}

// GetGolangImage 获取Golang镜像
func GetGolangImage(version string) string {
	return GetImageRegistry().GetImage("golang", version, true)
}
