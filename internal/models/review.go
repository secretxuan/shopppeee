package models

import (
	"time"

	"gorm.io/gorm"
)

// Review 商品评价模型
type Review struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID    uint   `gorm:"index;not null" json:"user_id"`
	ProductID uint   `gorm:"index;not null" json:"product_id"`
	OrderID   uint   `gorm:"index" json:"order_id"`
	Rating    int    `gorm:"not null" json:"rating" binding:"required,gte=1,lte=5"`
	Content   string `gorm:"type:text" json:"content"`
	Images    string `gorm:"type:text" json:"images"` // JSON数组字符串
	Status    string `gorm:"size:20;default:'pending'" json:"status"` // pending, approved, rejected

	// 关联
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 指定表名
func (Review) TableName() string {
	return "reviews"
}
