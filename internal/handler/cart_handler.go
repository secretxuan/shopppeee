package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shoppee/ecommerce/internal/service"
	"github.com/shoppee/ecommerce/pkg/response"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler() *CartHandler {
	return &CartHandler{
		cartService: service.NewCartService(),
	}
}

// GetCart 获取购物车
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取购物车失败")
		return
	}

	response.Success(c, cart)
}

// AddCartItem 添加商品到购物车
func (h *CartHandler) AddCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.cartService.AddCartItem(userID, req.ProductID, req.Quantity); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "添加成功"})
}

// UpdateCartItem 更新购物车项数量
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的购物车项ID")
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.cartService.UpdateCartItemQuantity(userID, uint(id), req.Quantity); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// DeleteCartItem 删除购物车项
func (h *CartHandler) DeleteCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的购物车项ID")
		return
	}

	if err := h.cartService.DeleteCartItem(userID, uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// SelectCartItem 切换购物车项选中状态
func (h *CartHandler) SelectCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的购物车项ID")
		return
	}

	var req struct {
		Selected bool `json:"selected"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.cartService.ToggleCartItemSelection(userID, uint(id), req.Selected); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "操作成功"})
}

// ClearCart 清空购物车
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.cartService.ClearCart(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "清空购物车失败")
		return
	}

	response.Success(c, gin.H{"message": "清空成功"})
}
