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

// OrderService 订单服务
type OrderService struct {
	cartService *CartService
}

// NewOrderService 创建订单服务实例
func NewOrderService() *OrderService {
	return &OrderService{
		cartService: NewCartService(),
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(userID uint, req interface{}) (*models.Order, error) {
	type OrderReq struct {
		AddressID     uint   `json:"address_id"`
		CartItemIDs   []uint `json:"cart_item_ids"`
		PaymentMethod string `json:"payment_method"`
		Remark        string `json:"remark"`
	}
	
	reqStruct, ok := req.(*OrderReq)
	if !ok {
		return nil, errors.New("无效的请求数据")
	}
	
	// 获取地址信息
	var address models.Address
	if err := database.DB.Where("id = ? AND user_id = ?", reqStruct.AddressID, userID).First(&address).Error; err != nil {
		return nil, errors.New("地址不存在")
	}
	
	// 获取购物车项
	var cartItems []models.CartItem
	if err := database.DB.Preload("Product").Where("id IN ?", reqStruct.CartItemIDs).Find(&cartItems).Error; err != nil {
		return nil, err
	}
	
	if len(cartItems) == 0 {
		return nil, errors.New("购物车为空")
	}
	
	// 计算总金额并创建订单
	var order *models.Order
	err := database.Transaction(func(tx *gorm.DB) error {
		// 计算总金额
		var totalAmount float64
		var orderItems []models.OrderItem
		
		for _, item := range cartItems {
			// 检查库存
			if item.Product.Stock < item.Quantity {
				return fmt.Errorf("商品 %s 库存不足", item.Product.Name)
			}
			
			subTotal := item.Product.Price * float64(item.Quantity)
			totalAmount += subTotal
			
			// 创建订单项
			orderItem := models.OrderItem{
				ProductID:    item.ProductID,
				Quantity:     item.Quantity,
				Price:        item.Product.Price,
				SubTotal:     subTotal,
				ProductName:  item.Product.Name,
				ProductSKU:   item.Product.SKU,
			}
			orderItems = append(orderItems, orderItem)
			
			// 扣减库存
			if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}
		
		// 生成订单号
		orderNo := fmt.Sprintf("ORD%d%d", time.Now().Unix(), userID)
		
		// 创建订单
		order = &models.Order{
			OrderNo:         orderNo,
			UserID:          userID,
			TotalAmount:     totalAmount,
			Status:          "pending",
			PaymentMethod:   reqStruct.PaymentMethod,
			PaymentStatus:   "unpaid",
			ReceiverName:    address.Name,
			ReceiverPhone:   address.Phone,
			ReceiverAddress: fmt.Sprintf("%s%s%s%s", address.Province, address.City, address.District, address.Detail),
			Remark:          reqStruct.Remark,
		}
		
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		
		// 创建订单项
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}
		
		// 删除购物车项
		if err := tx.Delete(&models.CartItem{}, reqStruct.CartItemIDs).Error; err != nil {
			return err
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	logger.Info("创建订单成功", zap.Uint("user_id", userID), zap.Uint("order_id", order.ID))
	return order, nil
}

// GetUserOrders 获取用户订单列表
func (s *OrderService) GetUserOrders(userID uint, page, pageSize int, status string) ([]models.Order, int64, error) {
	query := database.DB.Model(&models.Order{}).Where("user_id = ?", userID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	var orders []models.Order
	offset := (page - 1) * pageSize
	if err := query.Preload("OrderItems.Product").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	
	return orders, total, nil
}

// GetOrder 获取订单详情
func (s *OrderService) GetOrder(orderID, userID uint) (*models.Order, error) {
	var order models.Order
	if err := database.DB.Preload("OrderItems.Product").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(orderID, userID uint) error {
	return database.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
			return err
		}
		
		// 只能取消待支付的订单
		if order.Status != "pending" {
			return errors.New("订单状态不允许取消")
		}
		
		// 更新订单状态
		now := time.Now()
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":       "cancelled",
			"cancelled_at": &now,
		}).Error; err != nil {
			return err
		}
		
		// 恢复库存
		var orderItems []models.OrderItem
		if err := tx.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
			return err
		}
		
		for _, item := range orderItems {
			if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
				UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}
		
		logger.Info("取消订单成功", zap.Uint("order_id", orderID))
		return nil
	})
}

// ConfirmReceipt 确认收货
func (s *OrderService) ConfirmReceipt(orderID, userID uint) error {
	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return err
	}
	
	if order.Status != "shipped" {
		return errors.New("订单状态不正确")
	}
	
	now := time.Now()
	if err := database.DB.Model(&order).Updates(map[string]interface{}{
		"status":       "completed",
		"completed_at": &now,
	}).Error; err != nil {
		return err
	}
	
	logger.Info("确认收货成功", zap.Uint("order_id", orderID))
	return nil
}

// AdminGetOrders 管理员获取订单列表
func (s *OrderService) AdminGetOrders(page, pageSize int, status string) ([]models.Order, int64, error) {
	query := database.DB.Model(&models.Order{})
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	var orders []models.Order
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Preload("OrderItems.Product").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	
	return orders, total, nil
}

// AdminUpdateOrderStatus 管理员更新订单状态
func (s *OrderService) AdminUpdateOrderStatus(orderID uint, status string) error {
	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return err
	}
	
	updates := map[string]interface{}{"status": status}
	now := time.Now()
	
	switch status {
	case "paid":
		updates["payment_status"] = "paid"
		updates["paid_at"] = &now
	case "shipped":
		updates["shipped_at"] = &now
	case "completed":
		updates["completed_at"] = &now
	case "cancelled":
		updates["cancelled_at"] = &now
	}
	
	if err := database.DB.Model(&order).Updates(updates).Error; err != nil {
		return err
	}
	
	logger.Info("管理员更新订单状态", zap.Uint("order_id", orderID), zap.String("status", status))
	return nil
}
