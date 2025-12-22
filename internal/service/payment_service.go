package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/shoppee/ecommerce/internal/database"
	"github.com/shoppee/ecommerce/internal/models"
	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PaymentService 支付服务
type PaymentService struct{}

// NewPaymentService 创建支付服务实例
func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// CreatePayment 创建支付
func (s *PaymentService) CreatePayment(userID, orderID uint, paymentMethod string) (*models.Payment, error) {
	// 获取订单信息
	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return nil, errors.New("订单不存在")
	}
	
	if order.PaymentStatus == "paid" {
		return nil, errors.New("订单已支付")
	}
	
	if order.Status == "cancelled" {
		return nil, errors.New("订单已取消")
	}
	
	// 生成支付单号
	paymentNo := fmt.Sprintf("PAY%d%d", time.Now().Unix(), orderID)
	
	// 创建支付记录
	payment := &models.Payment{
		OrderID:       orderID,
		PaymentNo:     paymentNo,
		PaymentMethod: paymentMethod,
		Amount:        order.TotalAmount,
		Status:        "pending",
	}
	
	if err := database.DB.Create(payment).Error; err != nil {
		return nil, err
	}
	
	logger.Info("创建支付成功", zap.Uint("payment_id", payment.ID), zap.String("payment_no", paymentNo))
	
	// 这里应该调用第三方支付接口，这里简化处理
	// 返回支付信息给前端，前端跳转到支付页面
	
	return payment, nil
}

// GetPayment 获取支付详情
func (s *PaymentService) GetPayment(paymentID, userID uint) (*models.Payment, error) {
	var payment models.Payment
	if err := database.DB.Preload("Order", "user_id = ?", userID).
		First(&payment, paymentID).Error; err != nil {
		return nil, err
	}
	
	if payment.Order == nil {
		return nil, errors.New("无权访问该支付记录")
	}
	
	return &payment, nil
}

// HandlePaymentCallback 处理支付回调
func (s *PaymentService) HandlePaymentCallback(paymentNo, thirdPartyNo, status string) error {
	return database.Transaction(func(tx *gorm.DB) error {
		var payment models.Payment
		if err := tx.Where("payment_no = ?", paymentNo).First(&payment).Error; err != nil {
			return err
		}
		
		now := time.Now()
		
		if status == "success" {
			// 更新支付状态
			if err := tx.Model(&payment).Updates(map[string]interface{}{
				"status":         "success",
				"third_party_no": thirdPartyNo,
				"paid_at":        &now,
			}).Error; err != nil {
				return err
			}
			
			// 更新订单状态
			if err := tx.Model(&models.Order{}).Where("id = ?", payment.OrderID).Updates(map[string]interface{}{
				"payment_status": "paid",
				"status":         "paid",
				"paid_at":        &now,
			}).Error; err != nil {
				return err
			}
			
			logger.Info("支付成功", zap.String("payment_no", paymentNo))
		} else {
			// 支付失败
			if err := tx.Model(&payment).Update("status", "failed").Error; err != nil {
				return err
			}
			
			logger.Warn("支付失败", zap.String("payment_no", paymentNo))
		}
		
		return nil
	})
}
