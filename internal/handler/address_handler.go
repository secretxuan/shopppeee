package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shoppee/ecommerce/internal/service"
	"github.com/shoppee/ecommerce/pkg/response"
)

// AddressHandler 地址处理器
type AddressHandler struct {
	addressService *service.AddressService
}

// NewAddressHandler 创建地址处理器实例
func NewAddressHandler() *AddressHandler {
	return &AddressHandler{
		addressService: service.NewAddressService(),
	}
}

// GetAddressList 获取用户地址列表
func (h *AddressHandler) GetAddressList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	addresses, err := h.addressService.GetUserAddresses(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取地址列表失败")
		return
	}
	
	response.Success(c, addresses)
}

// GetAddress 获取地址详情
func (h *AddressHandler) GetAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的地址ID")
		return
	}
	
	address, err := h.addressService.GetAddress(uint(id), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "地址不存在")
		return
	}
	
	response.Success(c, address)
}

// CreateAddress 创建收货地址
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		Name      string `json:"name" binding:"required"`
		Phone     string `json:"phone" binding:"required"`
		Province  string `json:"province" binding:"required"`
		City      string `json:"city" binding:"required"`
		District  string `json:"district" binding:"required"`
		Detail    string `json:"detail" binding:"required"`
		IsDefault bool   `json:"is_default"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	
	address, err := h.addressService.CreateAddress(userID.(uint), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建地址失败")
		return
	}
	
	response.Success(c, address)
}

// UpdateAddress 更新收货地址
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的地址ID")
		return
	}
	
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	
	if err := h.addressService.UpdateAddress(uint(id), userID.(uint), req); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新地址失败")
		return
	}
	
	response.Success(c, gin.H{"message": "更新成功"})
}

// DeleteAddress 删除收货地址
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的地址ID")
		return
	}
	
	if err := h.addressService.DeleteAddress(uint(id), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除地址失败")
		return
	}
	
	response.Success(c, gin.H{"message": "删除成功"})
}

// SetDefaultAddress 设置默认地址
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的地址ID")
		return
	}
	
	if err := h.addressService.SetDefaultAddress(uint(id), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, "设置默认地址失败")
		return
	}
	
	response.Success(c, gin.H{"message": "设置成功"})
}

// GetDefaultAddress 获取默认地址
func (h *AddressHandler) GetDefaultAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	address, err := h.addressService.GetDefaultAddress(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "未找到默认地址")
		return
	}
	
	response.Success(c, address)
}
