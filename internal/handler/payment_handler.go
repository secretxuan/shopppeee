package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shoppee/ecommerce/internal/service"
	"github.com/shoppee/ecommerce/pkg/response"
)

// PaymentHandler 支付处理器
type PaymentHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentHandler 创建支付处理器实例
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService: service.NewPaymentService(),
	}
}

// CreatePayment 创建支付
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		OrderID       uint   `json:"order_id" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required,oneof=alipay wechat card"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	
	payment, err := h.paymentService.CreatePayment(userID.(uint), req.OrderID, req.PaymentMethod)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建支付失败: "+err.Error())
		return
	}
	
	response.Success(c, payment)
}

// GetPayment 获取支付详情
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的支付ID")
		return
	}
	
	payment, err := h.paymentService.GetPayment(uint(id), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "支付记录不存在")
		return
	}
	
	response.Success(c, payment)
}

// PaymentCallback 支付回调（模拟）
func (h *PaymentHandler) PaymentCallback(c *gin.Context) {
	var req struct {
		PaymentNo     string `json:"payment_no" binding:"required"`
		ThirdPartyNo  string `json:"third_party_no"`
		Status        string `json:"status" binding:"required,oneof=success failed"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	
	if err := h.paymentService.HandlePaymentCallback(req.PaymentNo, req.ThirdPartyNo, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, "处理支付回调失败: "+err.Error())
		return
	}
	
	response.Success(c, gin.H{"message": "支付回调处理成功"})
}
