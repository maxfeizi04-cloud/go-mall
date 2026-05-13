package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Log 全局日志实例
// 使用 SugaredLogger, 支持 printf 风格 (如 Log.Infof("xxx %s", name)

var Log *zap.SugaredLogger

// Init 初始化日志
// 参数:
//
//	     level		- 日志级别: debug/info/warn/error
//			filename	- 日志文件路径,如 "logs/app.log"
func Init(lever, filename string) {
	// ---- 1. 解析日志级别 ----
	var zapLever zapcore.Level
	switch lever {
	case "debug":
		zapLever = zapcore.DebugLevel // 输出所有日志
	case "info":
		zapLever = zapcore.InfoLevel // 输出 info 及以上
	case "warn":
		zapLever = zapcore.WarnLevel // 输出 warn 及以上
	case "error":
		zapLever = zapcore.ErrorLevel // 只输出 error
	default:
		zapLever = zapcore.InfoLevel
	}

	// ---- 2. 配置日志格式 ----
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                         // 时间字段在 JSON 中的 key
		LevelKey:       "level",                        // 日志级别字段的 key
		CallerKey:      "caller",                       // 调用位置（文件:行号）字段的 key
		MessageKey:     "msg",                          // 日志消息字段的 key
		LineEnding:     zapcore.DefaultLineEnding,      // 行结束符（默认 \n）
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 级别编码为小写
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // 时间编码格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 持续时间编码
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 调用者编码（短格式）
	}

	// ---- 3. 设置输出目标 ----

	// 控制台输出: 用可读的文本格式
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 文件输出: 用 JSON 格式 (方便后续用 ELK 等工具解析)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 文件写入器 (带日志轮转功能)
	fileWriter := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,  // 单位 MB
		MaxBackups: 5,    // 最多保留几个旧文件
		MaxAge:     30,   // 保留天数
		Compress:   true, // 压缩旧文件
	}

	// ---- 4. 创建日志核心 ----
	// Tee: 同时输出到多个目标 (控制台 + 文件)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLever),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), zapLever),
	)

	// ---- 5. 创建 Logger ----
	// AddCaller(): 记录调用日志的代码位置 (文件名:行号)
	// AddCallerSkip(1): 跳过一层调用栈,这样显示的是业务代码的位置而不是 logger 包的位置
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 转成 SugaredLogger 并赋值给全局变量
	Log = logger.Sugar()
}
