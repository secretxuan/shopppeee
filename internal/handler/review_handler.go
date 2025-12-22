package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shoppee/ecommerce/internal/service"
	"github.com/shoppee/ecommerce/pkg/response"
)

// ReviewHandler 评价处理器
type ReviewHandler struct {
	reviewService *service.ReviewService
}

// NewReviewHandler 创建评价处理器实例
func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		reviewService: service.NewReviewService(),
	}
}

// CreateReview 创建商品评价
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		ProductID uint   `json:"product_id" binding:"required"`
		OrderID   uint   `json:"order_id"`
		Rating    int    `json:"rating" binding:"required,gte=1,lte=5"`
		Content   string `json:"content"`
		Images    string `json:"images"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	
	review, err := h.reviewService.CreateReview(userID.(uint), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建评价失败: "+err.Error())
		return
	}
	
	response.Success(c, review)
}

// GetProductReviews 获取商品评价列表
func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的商品ID")
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	reviews, total, err := h.reviewService.GetProductReviews(uint(id), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取评价列表失败")
		return
	}
	
	response.SuccessWithPagination(c, reviews, total, page, pageSize)
}

// GetMyReviews 获取我的评价
func (h *ReviewHandler) GetMyReviews(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	reviews, total, err := h.reviewService.GetUserReviews(userID.(uint), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取评价列表失败")
		return
	}
	
	response.SuccessWithPagination(c, reviews, total, page, pageSize)
}

// DeleteReview 删除评价
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评价ID")
		return
	}
	
	if err := h.reviewService.DeleteReview(uint(id), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除评价失败: "+err.Error())
		return
	}
	
	response.Success(c, gin.H{"message": "删除成功"})
}

// AdminReplyReview 管理员回复评价
func (h *ReviewHandler) AdminReplyReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评价ID")
		return
	}
	
	var req struct {
		Reply string `json:"reply" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	
	if err := h.reviewService.ReplyReview(uint(id), req.Reply); err != nil {
		response.Error(c, http.StatusInternalServerError, "回复评价失败: "+err.Error())
		return
	}
	
	response.Success(c, gin.H{"message": "回复成功"})
}
