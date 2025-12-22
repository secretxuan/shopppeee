package service

import (
	"errors"
	"time"

	"github.com/shoppee/ecommerce/internal/database"
	"github.com/shoppee/ecommerce/internal/models"
	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

// ReviewService 评价服务
type ReviewService struct{}

// NewReviewService 创建评价服务实例
func NewReviewService() *ReviewService {
	return &ReviewService{}
}

// CreateReview 创建评价
func (s *ReviewService) CreateReview(userID uint, req interface{}) (*models.Review, error) {
	type ReviewReq struct {
		ProductID uint   `json:"product_id"`
		OrderID   uint   `json:"order_id"`
		Rating    int    `json:"rating"`
		Content   string `json:"content"`
		Images    string `json:"images"`
	}
	
	reqStruct, ok := req.(*ReviewReq)
	if !ok {
		return nil, errors.New("无效的请求数据")
	}
	
	// 检查是否已评价
	var existingReview models.Review
	if err := database.DB.Where("user_id = ? AND product_id = ? AND order_id = ?", 
		userID, reqStruct.ProductID, reqStruct.OrderID).First(&existingReview).Error; err == nil {
		return nil, errors.New("该订单商品已评价")
	}
	
	// 创建评价
	review := &models.Review{
		UserID:    userID,
		ProductID: reqStruct.ProductID,
		OrderID:   reqStruct.OrderID,
		Rating:    reqStruct.Rating,
		Content:   reqStruct.Content,
		Images:    reqStruct.Images,
		Status:    "published",
	}
	
	if err := database.DB.Create(review).Error; err != nil {
		return nil, err
	}
	
	logger.Info("创建评价成功", zap.Uint("user_id", userID), zap.Uint("review_id", review.ID))
	return review, nil
}

// GetProductReviews 获取商品评价列表
func (s *ReviewService) GetProductReviews(productID uint, page, pageSize int) ([]models.Review, int64, error) {
	query := database.DB.Model(&models.Review{}).
		Where("product_id = ? AND status = ?", productID, "published")
	
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	var reviews []models.Review
	offset := (page - 1) * pageSize
	if err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}
	
	return reviews, total, nil
}

// GetUserReviews 获取用户评价列表
func (s *ReviewService) GetUserReviews(userID uint, page, pageSize int) ([]models.Review, int64, error) {
	query := database.DB.Model(&models.Review{}).Where("user_id = ?", userID)
	
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	var reviews []models.Review
	offset := (page - 1) * pageSize
	if err := query.Preload("Product").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}
	
	return reviews, total, nil
}

// DeleteReview 删除评价
func (s *ReviewService) DeleteReview(reviewID, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", reviewID, userID).Delete(&models.Review{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("评价不存在")
	}
	
	logger.Info("删除评价成功", zap.Uint("review_id", reviewID))
	return nil
}

// ReplyReview 回复评价（管理员）
func (s *ReviewService) ReplyReview(reviewID uint, reply string) error {
	var review models.Review
	if err := database.DB.First(&review, reviewID).Error; err != nil {
		return err
	}
	
	now := time.Now()
	if err := database.DB.Model(&review).Updates(map[string]interface{}{
		"reply":      reply,
		"replied_at": &now,
	}).Error; err != nil {
		return err
	}
	
	logger.Info("回复评价成功", zap.Uint("review_id", reviewID))
	return nil
}
