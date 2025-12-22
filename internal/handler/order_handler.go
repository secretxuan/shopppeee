package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shoppee/ecommerce/internal/service"
	"github.com/shoppee/ecommerce/pkg/response"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler 创建订单处理器实例
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(),
	}
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		AddressID     uint   `json:"address_id" binding:"required"`
		CartItemIDs   []uint `json:"cart_item_ids" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required,oneof=alipay wechat card"`
		Remark        string `json:"remark"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	
	order, err := h.orderService.CreateOrder(userID.(uint), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建订单失败: "+err.Error())
		return
	}
	
	response.Success(c, order)
}

// GetOrderList 获取订单列表
func (h *OrderHandler) GetOrderList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	
	orders, total, err := h.orderService.GetUserOrders(userID.(uint), page, pageSize, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取订单列表失败")
		return
	}
	
	response.SuccessWithPagination(c, orders, total, page, pageSize)
}

// GetOrder 获取订单详情
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID")
		return
	}
	
	order, err := h.orderService.GetOrder(uint(id), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "订单不存在")
		return
	}
	
	response.Success(c, order)
}

// CancelOrder 取消订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID")
		return
	}
	
	if err := h.orderService.CancelOrder(uint(id), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, "取消订单失败: "+err.Error())
		return
	}
	
	response.Success(c, gin.H{"message": "订单已取消"})
}

// ConfirmReceipt 确认收货
func (h *OrderHandler) ConfirmReceipt(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID")
		return
	}
	
	if err := h.orderService.ConfirmReceipt(uint(id), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, "确认收货失败: "+err.Error())
		return
	}
	
	response.Success(c, gin.H{"message": "确认收货成功"})
}

// AdminGetOrderList 管理员获取订单列表
func (h *OrderHandler) AdminGetOrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	
	orders, total, err := h.orderService.AdminGetOrders(page, pageSize, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取订单列表失败")
		return
	}
	
	response.SuccessWithPagination(c, orders, total, page, pageSize)
}

// AdminUpdateOrderStatus 管理员更新订单状态
func (h *OrderHandler) AdminUpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID")
		return
	}
	
	var req struct {
		Status string `json:"status" binding:"required,oneof=pending paid shipped completed cancelled"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	
	if err := h.orderService.AdminUpdateOrderStatus(uint(id), req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新订单状态失败: "+err.Error())
		return
	}
	
	response.Success(c, gin.H{"message": "订单状态已更新"})
}
