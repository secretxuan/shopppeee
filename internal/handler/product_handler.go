package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shoppee/ecommerce/internal/service"
	"github.com/shoppee/ecommerce/pkg/response"
)

// ProductHandler 商品处理器
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler 创建商品处理器实例
func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: service.NewProductService(),
	}
}

// GetProductList 获取商品列表
// @Summary 获取商品列表
// @Description 分页获取商品列表，支持筛选和排序
// @Tags 商品
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category_id query int false "分类ID"
// @Param keyword query string false "搜索关键词"
// @Param sort query string false "排序方式" Enums(price_asc, price_desc, sale_desc, new)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /products [get]
func (h *ProductHandler) GetProductList(c *gin.Context) {
	var req service.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	products, total, err := h.productService.GetProductList(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取商品列表失败")
		return
	}

	response.Page(c, products, total, req.Page, req.PageSize)
}

// GetProductByID 获取商品详情
// @Summary 获取商品详情
// @Description 根据ID获取商品详细信息
// @Tags 商品
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response{data=models.Product}
// @Failure 404 {object} response.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的商品ID")
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "商品不存在")
		return
	}

	response.Success(c, product)
}

// SearchProducts 搜索商品
// @Summary 搜索商品
// @Description 根据关键词搜索商品
// @Tags 商品
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	products, total, err := h.productService.SearchProducts(keyword, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Page(c, products, total, page, pageSize)
}

// BatchUpdateStock 批量更新库存
// @Summary 批量更新库存
// @Description 批量更新商品库存（需要管理员权限）
// @Tags 商品
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param updates body map[uint]int true "库存更新映射"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /products/batch-stock [post]
func (h *ProductHandler) BatchUpdateStock(c *gin.Context) {
	var updates map[uint]int
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.productService.BatchUpdateStock(updates); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "批量更新库存成功", nil)
}
