package models

import (
	"time"

	"gorm.io/gorm"
)

// Order 订单模型
type Order struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OrderNo     string  `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	UserID      uint    `gorm:"index;not null" json:"user_id"`
	TotalAmount float64 `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	PayAmount   float64 `gorm:"type:decimal(10,2);not null" json:"pay_amount"`
	Status      string  `gorm:"size:20;default:'pending';index" json:"status"` // pending, paid, shipped, completed, cancelled
	PayStatus   string  `gorm:"size:20;default:'unpaid'" json:"pay_status"`    // unpaid, paid, refunded
	PayMethod   string  `gorm:"size:20" json:"pay_method"`                     // alipay, wechat, credit_card
	Remark      string  `gorm:"type:text" json:"remark"`

	// 收货信息
	ReceiverName    string `gorm:"size:50" json:"receiver_name"`
	ReceiverPhone   string `gorm:"size:20" json:"receiver_phone"`
	ReceiverAddress string `gorm:"size:255" json:"receiver_address"`

	// 时间字段
	PaidAt      *time.Time `json:"paid_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CancelledAt *time.Time `json:"cancelled_at"`

	// 关联
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
	Payment    *Payment    `gorm:"foreignKey:OrderID" json:"payment,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项模型
type OrderItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OrderID   uint    `gorm:"index;not null" json:"order_id"`
	ProductID uint    `gorm:"index;not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity" binding:"required,gt=0"`
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`

	// 快照字段（防止商品信息变更）
	ProductName  string `gorm:"size:200" json:"product_name"`
	ProductImage string `gorm:"size:255" json:"product_image"`
	ProductSKU   string `gorm:"size:100" json:"product_sku"`

	// 关联
	Order   *Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}
