package models

import (
	"time"

	"gorm.io/gorm"
)

// Product 商品模型
type Product struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string  `gorm:"size:200;not null;index" json:"name" binding:"required"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price" binding:"required,gt=0"`
	OrigPrice   float64 `gorm:"type:decimal(10,2)" json:"orig_price"` // 原价
	Stock       int     `gorm:"not null;default:0" json:"stock" binding:"gte=0"`
	SKU         string  `gorm:"uniqueIndex;size:100" json:"sku"`
	Images      string  `gorm:"type:text" json:"images"` // JSON数组字符串
	Status      string  `gorm:"size:20;default:'active'" json:"status"` // active, inactive, out_of_stock
	ViewCount   int     `gorm:"default:0" json:"view_count"`
	SaleCount   int     `gorm:"default:0" json:"sale_count"`

	// 外键
	CategoryID uint      `gorm:"index" json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	// 关联
	Reviews   []Review    `gorm:"foreignKey:ProductID" json:"reviews,omitempty"`
	CartItems []CartItem  `gorm:"foreignKey:ProductID" json:"-"`
	OrderItems []OrderItem `gorm:"foreignKey:ProductID" json:"-"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// Category 商品分类模型
type Category struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"size:100;not null;index" json:"name" binding:"required"`
	Description string `gorm:"type:text" json:"description"`
	Icon        string `gorm:"size:255" json:"icon"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      string `gorm:"size:20;default:'active'" json:"status"`

	// 父子分类
	ParentID *uint      `gorm:"index" json:"parent_id"`
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	// 关联
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}
