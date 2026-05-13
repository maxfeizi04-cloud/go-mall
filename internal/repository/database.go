package repository

import (
	"fmt"
	"time"

	"github.com/maxfeizi04-cloud/go-mall/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
// 所有 repository 操作都通过这个变量访问数据库
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg config.DatabaseConfig) {
	var err error

	// 打开数据库连接
	// cfg.DSN() 返回 "root:root123@tcp(localhost:3306)/go-mall?..." 格式的连接字符串
	DB, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		// 设置 SQL 日志级别: Info 表示打印所有 SQL
		// 生成环境建议改成 logger.Warn (只打印慢查询和错误)
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}

	// 获取底层的 *sql.DB 对象,用于配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("获取 DB 实例失败: %v", err))
	}

	// ---- 连接池配置 ----
	// MaxIdleConns: 空闲连接数. 设太大浪费内存,设太小频繁创建连接
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// MaxOpenConns: 最大打开连接数.根据数据库承受能力设置
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// ConnMaxLifetime: 连接最大存活时间.设置为 1 小时,防止使用过期连接
	sqlDB.SetConnMaxLifetime(time.Hour)
}
