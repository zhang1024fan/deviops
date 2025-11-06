// 文件配置,解析yaml配种文件
// author xiaoRui

package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// 总配文件
type config struct {
	Server        server        `yaml:"server"`
	Db            db            `yaml:"db"`
	Redis         redis         `yaml:"redis"`
	ImageSettings imageSettings `yaml:"imageSettings"`
	Log           log           `yaml:"log"`
	Monitor       monitor       `yaml:"monitor"`
}

// 监控配置
type monitor struct {
	Prometheus  prometheus  `yaml:"prometheus"`
	Pushgateway pushgateway `yaml:"pushgateway"`
	Agent       agent       `yaml:"agent"`
}

// Pushgateway配置
type pushgateway struct {
	URL string `yaml:"url"`
}

// Agent配置
type agent struct {
	HeartbeatServerURL string `yaml:"heartbeat_server_url"`
	HeartbeatToken     string `yaml:"heartbeat_token"`
}

// Prometheus配置
type prometheus struct {
	URL string `yaml:"url"`
}

// 项目端口配置
type server struct {
	Address       string `yaml:"address"`
	Model         string `yaml:"model"`
	EnableSwagger bool   `yaml:"enableSwagger"` // 是否启用Swagger文档
}

// 数据库配置
type db struct {
	Dialects string `yaml:"dialects"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Db       string `yaml:"db"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
	MaxIdle  int    `yaml:"maxIdle"`
	MaxOpen  int    `yaml:"maxOpen"`
}

// redis配置
type redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

// imageSettings图片上传配置
type imageSettings struct {
	UploadDir string `yaml:"uploadDir"`
	ImageHost string `yaml:"imageHost"`
}

// log日志配置
type log struct {
	Path  string `yaml:"path"`
	Name  string `yaml:"name"`
	Model string `yaml:"model"`
}

var Config *config

// 配置初始化
func init() {
	// 初始化时先不加载配置文件，等待LoadConfig()被调用
}

// LoadConfig 从指定路径加载配置文件
func LoadConfig(configPath string) error {
	if configPath == "" {
		configPath = "./config.yaml" // 默认配置文件路径
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// 绑定值
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		return err
	}

	return nil
}

// GetConfig 获取数据库配置
func GetConfig() *db {
	if Config == nil {
		panic("Config is not initialized")
	}
	return &Config.Db
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig() *redis {
	if Config == nil {
		panic("Config is not initialized")
	}
	return &Config.Redis
}

// Setup 初始化配置（为了兼容migrate.go的调用）
func Setup() {
	// 配置已经在init()方法中初始化了，这里只是提供一个兼容性方法
	if Config == nil {
		panic("Config initialization failed")
	}
}
