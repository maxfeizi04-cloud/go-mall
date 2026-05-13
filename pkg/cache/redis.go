package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RDB 全局 Redis 客户端实例
var RDB *redis.Client

// RedisConfig Redis 配置结构体
// 为什么单独定义而不是直接用 config.RedisConfig？
// 因为 pkg 包不应该依赖 internal 包（Go 的包可见性规则）
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// Init 初始化 Redis 连接
func Init(cfg RedisConfig) {
	// 创建 Redis 客户端
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // 地址: host:port
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 用 Ping 测试连接是否成功
	// 设置 5s 超时,防止卡死
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RDB.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis 连接失败: %s", err))
	}

}

// Set 设置缓存
// 参数 expiration 是过期时间, 0 表示永不过期
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RDB.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
// 返回值是字符串，如果 key 不存在会返回 redis.Nil 错误
func Get(ctx context.Context, key string) (interface{}, error) {
	return RDB.Get(ctx, key).Result()
}

// Del 删除缓存 (支持一次删除多个 key)
func Del(ctx context.Context, key ...string) error {
	return RDB.Del(ctx, key...).Err()
}

// Incr 自增(用于限流计数器等场景)
// 如果 key 不存在,会先初始化为 0 再加 1
func Incr(ctx context.Context, key string) (int64, error) {
	return RDB.Incr(ctx, key).Result()
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return RDB.Expire(ctx, key, expiration).Err()
}
