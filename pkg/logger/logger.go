package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// InitLogger 初始化日志系统
func InitLogger(level, logFilePath string) error {
	// 解析日志级别
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 创建日志目录
	if logFilePath != "" {
		dir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 控制台输出
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	var cores []zapcore.Core
	cores = append(cores, consoleCore)

	// 文件输出
	if logFilePath != "" {
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
			fileCore := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(file),
				zapLevel,
			)
			cores = append(cores, fileCore)
		}
	}

	// 组合多个Core
	core := zapcore.NewTee(cores...)

	// 创建Logger
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// Debug 输出调试日志
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info 输出信息日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 输出警告日志
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error 输出错误日志
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal 输出致命错误日志并退出
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Sync 刷新日志缓冲
func Sync() {
	_ = log.Sync()
}
