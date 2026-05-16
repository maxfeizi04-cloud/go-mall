package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxfeizi04-cloud/go-mall/internal/config"
	"github.com/maxfeizi04-cloud/go-mall/internal/model"
	"github.com/maxfeizi04-cloud/go-mall/internal/repository"
	"github.com/maxfeizi04-cloud/go-mall/pkg/cache"
	"github.com/maxfeizi04-cloud/go-mall/pkg/logger"
)

func main() {
	// ========== 1. 加载配置 ==========
	// 读取项目跟目录下的 config.yaml
	cfg, err := config.Load("config.yaml")
	if err != nil {
		// 配置加载失败是致命的错误,直接退出
		log.Fatalf("加载配置失败: %v", err)
	}

	// ========== 2. 初始化日志 ==========
	// 在所有组件之前初始化,这样后续组件的日志才能正常输出
	logger.Init(cfg.Log.Level, cfg.Log.FileName)
	logger.Log.Info("日志初始化完成")

	// ========== 3. 初始化数据库 ==========
	repository.InitDB(cfg.Database)
	logger.Log.Info("数据库连接成功")

	// ============================================================
	// 第 4 步：数据库自动迁移（自动建表）
	// GORM 的 AutoMigrate 会根据 struct 定义自动执行：
	//   - 建表（如果表不存在）
	//   - 添加新字段（如果字段不存在）
	//   - 添加索引（如果索引不存在）
	// 注意：AutoMigrate 不会删除已有字段或修改字段类型
	//       这是安全的，不会丢数据
	//
	// 传入的顺序建议按依赖关系：
	//   先建被依赖的表（User、Category），
	//   再建有外键的表（Product、CartItem、Order、OrderItem）
	// ============================================================
	err = repository.DB.AutoMigrate(
		&model.User{},      // → users 表
		&model.Category{},  // → categories 表
		&model.Product{},   // → products 表
		&model.CartItem{},  // → cart_items 表
		&model.Order{},     // → orders 表
		&model.OrderItem{}, // → order_items 表
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	logger.Log.Info("数据库迁移成功,6 张表已就绪")

	// ========== 5. 初始化 Redis ==========
	cache.Init(cache.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	logger.Log.Info("Redis 连接成功")

	// ========== 6. 初始化 Gin 并注册路由 ==========
	// 设置 Gin 的运行模式
	gin.SetMode(cfg.Sever.Mode)

	// 创建 Gin 引擎 (不带默认中间件)
	r := gin.Default()

	// Recovery 中间件: 捕获 panic,防止一个请求的 panic 导致整个服务崩溃
	r.Use(gin.Recovery())

	// 健康检查接口: 用于 Docker、负载均衡器探测服务是否存活
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// ========== 7. 启动 HTTP 服务 ==========
	addr := fmt.Sprintf(":%d", cfg.Sever.HTTPPort)
	logger.Log.Infof("HTTP 服务启动,监听地址: %s", addr)

	// r.Run 会阻塞,持续监听端口
	if err = r.Run(addr); err != nil {
		log.Fatalf("HTTP 服务启动失败: %v", err)
	}

}
