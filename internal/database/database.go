package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shoppee/ecommerce/internal/config"
	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	// DB 全局数据库实例
	DB *gorm.DB
	// RedisClient 全局Redis客户端
	RedisClient *redis.Client
)

// InitDB 初始化数据库连接
func InitDB() error {
	dsn := config.AppConfig.GetDSN()

	// 设置GORM日志级别
	logLevel := gormLogger.Silent
	if config.AppConfig.Debug {
		logLevel = gormLogger.Info
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取通用数据库对象sql.DB，设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	DB = db
	logger.Info("数据库连接成功")
	return nil
}

// InitRedis 初始化Redis连接
func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.GetRedisAddr(),
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
		PoolSize: 10,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	logger.Info("Redis连接成功")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// Transaction 事务处理辅助函数
func Transaction(fn func(*gorm.DB) error) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := fn(tx); err != nil {
			logger.Error("事务执行失败", zap.Error(err))
			return err
		}
		return nil
	})
}
