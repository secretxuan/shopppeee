package database

import (
	"github.com/shoppee/ecommerce/internal/models"
	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	logger.Info("开始数据库迁移...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Category{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
		&models.Address{},
		&models.Payment{},
		&models.Review{},
	)

	if err != nil {
		logger.Error("数据库迁移失败", zap.Error(err))
		return err
	}

	logger.Info("数据库迁移完成")
	return nil
}
