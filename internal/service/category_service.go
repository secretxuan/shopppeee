package service

import (
	"errors"

	"github.com/shoppee/ecommerce/internal/database"
	"github.com/shoppee/ecommerce/internal/models"
)

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// GetCategoryTree 获取分类树
func (s *CategoryService) GetCategoryTree() ([]*models.Category, error) {
	var categories []*models.Category
	
	// 获取所有分类
	if err := database.DB.Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	// 构建树形结构
	return buildTree(categories, 0), nil
}

// buildTree 构建树形结构
func buildTree(categories []*models.Category, parentID uint) []*models.Category {
	var tree []*models.Category
	
	for _, category := range categories {
		if category.ParentID == parentID {
			children := buildTree(categories, category.ID)
			if len(children) > 0 {
				category.Children = children
			}
			tree = append(tree, category)
		}
	}
	
	return tree
}

// GetCategoryList 获取分类列表（根据父ID）
func (s *CategoryService) GetCategoryList(parentID uint) ([]*models.Category, error) {
	var categories []*models.Category
	
	if err := database.DB.Where("parent_id = ?", parentID).
		Order("sort ASC, id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(name, description string, parentID uint, icon string, sort int, status string) (*models.Category, error) {
	// 验证父分类是否存在
	if parentID > 0 {
		var parent models.Category
		if err := database.DB.First(&parent, parentID).Error; err != nil {
			return nil, errors.New("父分类不存在")
		}
	}

	if status == "" {
		status = "active"
	}

	category := &models.Category{
		Name:        name,
		Description: description,
		ParentID:    parentID,
		Icon:        icon,
		Sort:        sort,
		Status:      status,
	}

	if err := database.DB.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(id uint, updates map[string]interface{}) error {
	// 检查分类是否存在
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		return errors.New("分类不存在")
	}

	return database.DB.Model(&category).Updates(updates).Error
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(id uint) error {
	// 检查是否有子分类
	var count int64
	database.DB.Model(&models.Category{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该分类下有子分类，无法删除")
	}

	// 检查是否有商品
	database.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该分类下有商品，无法删除")
	}

	return database.DB.Delete(&models.Category{}, id).Error
}
