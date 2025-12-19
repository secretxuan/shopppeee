package models

import (
	"time"

	"gorm.io/gorm"
)

// Payment 支付记录模型
type Payment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OrderID       uint    `gorm:"uniqueIndex;not null" json:"order_id"`
	TransactionNo string  `gorm:"size:100" json:"transaction_no"` // 第三方交易号
	PayMethod     string  `gorm:"size:20;not null" json:"pay_method"`
	PayAmount     float64 `gorm:"type:decimal(10,2);not null" json:"pay_amount"`
	Status        string  `gorm:"size:20;default:'pending'" json:"status"` // pending, success, failed, refunded
	PaidAt        *time.Time `json:"paid_at"`
	RefundedAt    *time.Time `json:"refunded_at"`

	// 关联
	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}
