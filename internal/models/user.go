package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username  string `gorm:"uniqueIndex;size:50;not null" json:"username" binding:"required,min=3,max=50"`
	Email     string `gorm:"uniqueIndex;size:100;not null" json:"email" binding:"required,email"`
	Password  string `gorm:"size:255;not null" json:"-"`
	Phone     string `gorm:"size:20" json:"phone"`
	Avatar    string `gorm:"size:255" json:"avatar"`
	Role      string `gorm:"size:20;default:'user'" json:"role"` // user, admin
	Status    string `gorm:"size:20;default:'active'" json:"status"` // active, inactive, banned
	LastLogin *time.Time `json:"last_login"`

	// 关联
	Addresses []Address `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	Orders    []Order   `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Cart      *Cart     `gorm:"foreignKey:UserID" json:"cart,omitempty"`
	Reviews   []Review  `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM钩子：创建前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// SetPassword 设置密码（加密）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
