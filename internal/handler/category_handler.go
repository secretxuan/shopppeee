package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shoppee/ecommerce/internal/service"
	"github.com/shoppee/ecommerce/pkg/response"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: service.NewCategoryService(),
	}
}

// GetCategoryTree 获取分类树
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	categories, err := h.categoryService.GetCategoryTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取分类树失败")
		return
	}

	response.Success(c, gin.H{
		"categories": categories,
	})
}

// GetCategoryList 获取分类列表
func (h *CategoryHandler) GetCategoryList(c *gin.Context) {
	parentIDStr := c.DefaultQuery("parent_id", "0")
	parentID, _ := strconv.ParseUint(parentIDStr, 10, 64)

	categories, err := h.categoryService.GetCategoryList(uint(parentID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取分类列表失败")
		return
	}

	response.Success(c, gin.H{
		"categories": categories,
	})
}

// CreateCategory 创建分类（管理员）
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		ParentID    uint   `json:"parent_id"`
		Icon        string `json:"icon"`
		Sort        int    `json:"sort"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	category, err := h.categoryService.CreateCategory(
		req.Name,
		req.Description,
		req.ParentID,
		req.Icon,
		req.Sort,
		req.Status,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建分类失败: "+err.Error())
		return
	}

	response.Success(c, category)
}

// UpdateCategory 更新分类（管理员）
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Sort        int    `json:"sort"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.categoryService.UpdateCategory(uint(id), updates); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新分类失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// DeleteCategory 删除分类（管理员）
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除分类失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
