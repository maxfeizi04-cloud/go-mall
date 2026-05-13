package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 应用总配置
// 每个字段的 mapstructure tag 对应 YAML 中的 key
// 比如 YAML 中的 "server.http_port" 对应 Server.HTTPPort
type Config struct {
	Sever    ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	HTTPPort int    `mapstructure:"http_port"`
	GRPCPort int    `mapstructure:"grpc_port"`
	Mode     string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// DSN 生成数据库连接字符串
// DSN = Data Source Name,是 MySQL 驱动要求的连接格式
// 格式: 用户名:密码@tcp(主机:端口)/数据库名?参数
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName,
	)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// LogConfig Log 配置
type LogConfig struct {
	Level    string `mapstructure:"level"`
	FileName string `mapstructure:"filename"`
}

// Load 加载配置文件
// 参数 path 是配置文件路径,比如 "config.yaml"
// 返回 Config 指针和错误
func Load(path string) (*Config, error) {
	// 告诉 viper 配置文件在哪
	viper.SetConfigFile(path)

	// 允许用环境变量覆盖配置(比如 DATABASE_HOST 环境变量会覆盖 database.host)
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 把读到的配置映射到 Config 结构体
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}
