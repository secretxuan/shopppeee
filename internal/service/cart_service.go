package service

import (
	"errors"

	"github.com/shoppee/ecommerce/internal/database"
	"github.com/shoppee/ecommerce/internal/models"
	"gorm.io/gorm"
)

type CartService struct{}

func NewCartService() *CartService {
	return &CartService{}
}

// GetCart 获取用户购物车
func (s *CartService) GetCart(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := database.DB.Preload("Items.Product").
		Where("user_id = ?", userID).
		First(&cart).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果购物车不存在，创建一个
			cart = models.Cart{UserID: userID}
			if err := database.DB.Create(&cart).Error; err != nil {
				return nil, err
			}
			return &cart, nil
		}
		return nil, err
	}

	return &cart, nil
}

// AddCartItem 添加商品到购物车
func (s *CartService) AddCartItem(userID, productID uint, quantity int) error {
	// 检查商品是否存在且有足够库存
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		return errors.New("商品不存在")
	}

	if product.Stock < quantity {
		return errors.New("库存不足")
	}

	// 获取或创建购物车
	cart, err := s.GetCart(userID)
	if err != nil {
		return err
	}

	// 检查是否已存在该商品
	var existingItem models.CartItem
	err = database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, productID).
		First(&existingItem).Error

	if err == nil {
		// 已存在，更新数量
		newQuantity := existingItem.Quantity + quantity
		if product.Stock < newQuantity {
			return errors.New("库存不足")
		}
		return database.DB.Model(&existingItem).Update("quantity", newQuantity).Error
	}

	// 不存在，创建新项
	item := models.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
		Selected:  true,
	}

	return database.DB.Create(&item).Error
}

// UpdateCartItemQuantity 更新购物车项数量
func (s *CartService) UpdateCartItemQuantity(userID, itemID uint, quantity int) error {
	// 获取购物车项
	var item models.CartItem
	err := database.DB.Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		Preload("Product").
		First(&item).Error
	
	if err != nil {
		return errors.New("购物车项不存在")
	}

	// 检查库存
	if item.Product.Stock < quantity {
		return errors.New("库存不足")
	}

	return database.DB.Model(&item).Update("quantity", quantity).Error
}

// DeleteCartItem 删除购物车项
func (s *CartService) DeleteCartItem(userID, itemID uint) error {
	result := database.DB.Where("id = ? AND cart_id IN (SELECT id FROM carts WHERE user_id = ?)", itemID, userID).
		Delete(&models.CartItem{})
	
	if result.RowsAffected == 0 {
		return errors.New("购物车项不存在")
	}

	return result.Error
}

// ToggleCartItemSelection 切换购物车项选中状态
func (s *CartService) ToggleCartItemSelection(userID, itemID uint, selected bool) error {
	result := database.DB.Model(&models.CartItem{}).
		Where("id = ? AND cart_id IN (SELECT id FROM carts WHERE user_id = ?)", itemID, userID).
		Update("selected", selected)
	
	if result.RowsAffected == 0 {
		return errors.New("购物车项不存在")
	}

	return result.Error
}

// ClearCart 清空购物车
func (s *CartService) ClearCart(userID uint) error {
	return database.DB.Where("cart_id IN (SELECT id FROM carts WHERE user_id = ?)", userID).
		Delete(&models.CartItem{}).Error
}

// GetSelectedItems 获取选中的购物车项
func (s *CartService) GetSelectedItems(userID uint) ([]models.CartItem, error) {
	var items []models.CartItem
	err := database.DB.Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ? AND cart_items.selected = ?", userID, true).
		Preload("Product").
		Find(&items).Error
	
	return items, err
}
