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

	UserID    uint   `gorm:"index;not null" json:"user_id"`
	Name      string `gorm:"size:50;not null" json:"name" binding:"required"`
	Phone     string `gorm:"size:20;not null" json:"phone" binding:"required"`
	Province  string `gorm:"size:50;not null" json:"province" binding:"required"`
	City      string `gorm:"size:50;not null" json:"city" binding:"required"`
	District  string `gorm:"size:50;not null" json:"district" binding:"required"`
	Detail    string `gorm:"size:255;not null" json:"detail" binding:"required"`
	IsDefault bool   `gorm:"default:false" json:"is_default"`
	
	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Address) TableName() string {
	return "addresses"
}

// BeforeCreate GORM钩子：创建前设置默认地址
func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.IsDefault {
		// 将该用户的其他地址设为非默认
		tx.Model(&Address{}).Where("user_id = ? AND is_default = ?", a.UserID, true).Update("is_default", false)
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前设置默认地址
func (a *Address) BeforeUpdate(tx *gorm.DB) error {
	if a.IsDefault {
		// 将该用户的其他地址设为非默认
		tx.Model(&Address{}).Where("user_id = ? AND id != ? AND is_default = ?", a.UserID, a.ID, true).Update("is_default", false)
	}
	return nil
}
