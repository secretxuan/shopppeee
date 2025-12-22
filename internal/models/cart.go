package models

import (
	"time"

	"gorm.io/gorm"
)

// Cart 购物车模型
type Cart struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`
	
	// 关联
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CartItems []CartItem `gorm:"foreignKey:CartID" json:"cart_items,omitempty"`
}

// TableName 指定表名
func (Cart) TableName() string {
	return "carts"
}

// CartItem 购物车项模型
type CartItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CartID    uint `gorm:"index;not null" json:"cart_id"`
	ProductID uint `gorm:"index;not null" json:"product_id"`
	Quantity  int  `gorm:"not null;default:1" json:"quantity"`
	Selected  bool `gorm:"default:true" json:"selected"` // 是否选中（用于结算）
	
	// 关联
	Cart    *Cart    `gorm:"foreignKey:CartID" json:"cart,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 指定表名
func (CartItem) TableName() string {
	return "cart_items"
}
