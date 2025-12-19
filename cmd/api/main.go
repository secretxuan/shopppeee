package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shoppee/ecommerce/internal/config"
	"github.com/shoppee/ecommerce/internal/database"
	"github.com/shoppee/ecommerce/internal/router"
	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

// @title Shoppee E-Commerce API
// @version 1.0
// @description 高并发电商系统API，基于Go + Gin框架
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.shoppee.com/support
// @contact.email support@shoppee.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT认证，格式: Bearer {token}

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		panic(fmt.Sprintf("初始化配置失败: %v", err))
	}

	// 初始化日志
	if err := logger.InitLogger(config.AppConfig.LogLevel, config.AppConfig.LogFilePath); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer logger.Sync()

	logger.Info("启动 Shoppee 电商系统", zap.String("env", config.AppConfig.Env))

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		logger.Fatal("数据库初始化失败", zap.Error(err))
	}

	// 自动迁移数据库
	if err := database.AutoMigrate(); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}

	// 初始化Redis
	if err := database.InitRedis(); err != nil {
		logger.Fatal("Redis初始化失败", zap.Error(err))
	}

	// 初始化路由
	r := router.SetupRouter()

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 启动服务器（优雅关闭）
	go func() {
		logger.Info("服务器启动", zap.Int("port", config.AppConfig.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 设置5秒超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("服务器强制关闭", zap.Error(err))
	}

	logger.Info("服务器已退出")
}
