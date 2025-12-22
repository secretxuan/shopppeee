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
	PaymentNo     string  `gorm:"uniqueIndex;size:50;not null" json:"payment_no"`
	PaymentMethod string  `gorm:"size:20;not null" json:"payment_method"` // alipay, wechat, card
	Amount        float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        string  `gorm:"size:20;default:'pending'" json:"status"` // pending, success, failed, refunded
	PaidAt        *time.Time `json:"paid_at"`
	RefundedAt    *time.Time `json:"refunded_at"`
	
	// 第三方支付信息
	ThirdPartyNo string `gorm:"size:100" json:"third_party_no"` // 第三方交易号
	
	// 备注
	Remark string `gorm:"type:text" json:"remark"`
	
	// 关联
	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}
