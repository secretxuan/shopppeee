package models

import (
	"time"

	"gorm.io/gorm"
)

// Address 收货地址模型
type Address struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID       uint   `gorm:"index;not null" json:"user_id"`
	ReceiverName string `gorm:"size:50;not null" json:"receiver_name" binding:"required"`
	Phone        string `gorm:"size:20;not null" json:"phone" binding:"required"`
	Province     string `gorm:"size:50" json:"province"`
	City         string `gorm:"size:50" json:"city"`
	District     string `gorm:"size:50" json:"district"`
	Detail       string `gorm:"size:255;not null" json:"detail" binding:"required"`
	IsDefault    bool   `gorm:"default:false" json:"is_default"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Address) TableName() string {
	return "addresses"
}
