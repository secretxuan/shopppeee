package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/shoppee/ecommerce/internal/database"
	"github.com/shoppee/ecommerce/internal/models"
	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProductService 商品服务
type ProductService struct{}

// NewProductService 创建商品服务实例
func NewProductService() *ProductService {
	return &ProductService{}
}

// ProductListRequest 商品列表请求
type ProductListRequest struct {
	Page       int    `form:"page" binding:"omitempty,gte=1"`
	PageSize   int    `form:"page_size" binding:"omitempty,gte=1,lte=100"`
	CategoryID uint   `form:"category_id"`
	Keyword    string `form:"keyword"`
	Sort       string `form:"sort"` // price_asc, price_desc, sale_desc, new
	Status     string `form:"status"`
}

// GetProductList 获取商品列表（支持分页、筛选、排序）
func (s *ProductService) GetProductList(req *ProductListRequest) ([]models.Product, int64, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	query := database.DB.Model(&models.Product{}).Preload("Category")

	// 分类筛选
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	} else {
		query = query.Where("status = ?", "active")
	}

	// 排序
	switch req.Sort {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "sale_desc":
		query = query.Order("sale_count DESC")
	case "new":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("id DESC")
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var products []models.Product
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetProductByID 根据ID获取商品详情（带缓存）
func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product

	// 先从Redis缓存获取
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%d", id)

	// 尝试从缓存读取（这里简化处理，实际应序列化完整对象）
	_, err := database.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中，从数据库加载完整信息
		if err := database.DB.Preload("Category").Preload("Reviews").First(&product, id).Error; err != nil {
			return nil, err
		}

		// 异步增加浏览量
		go s.incrementViewCount(id)

		return &product, nil
	}

	// 缓存未命中，从数据库查询
	if err := database.DB.Preload("Category").Preload("Reviews").First(&product, id).Error; err != nil {
		return nil, err
	}

	// 异步缓存到Redis（1小时过期）
	go func() {
		database.RedisClient.Set(ctx, cacheKey, "cached", 1*time.Hour)
	}()

	// 异步增加浏览量
	go s.incrementViewCount(id)

	return &product, nil
}

// incrementViewCount 增加商品浏览量
func (s *ProductService) incrementViewCount(productID uint) {
	database.DB.Model(&models.Product{}).Where("id = ?", productID).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
}

// BatchUpdateStock 批量更新库存（利用Go协程+Channel实现高并发）
func (s *ProductService) BatchUpdateStock(updates map[uint]int) error {
	if len(updates) == 0 {
		return errors.New("更新列表为空")
	}

	// 创建协程池处理批量更新
	const workerCount = 10
	jobs := make(chan struct {
		productID uint
		quantity  int
	}, len(updates))
	results := make(chan error, len(updates))

	// 启动worker协程池
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				err := s.updateStock(job.productID, job.quantity)
				results <- err
			}
		}()
	}

	// 发送任务到jobs channel
	for productID, quantity := range updates {
		jobs <- struct {
			productID uint
			quantity  int
		}{productID, quantity}
	}
	close(jobs)

	// 等待所有worker完成
	wg.Wait()
	close(results)

	// 检查是否有错误
	var errs []error
	for err := range results {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		logger.Error("批量更新库存失败", zap.Int("error_count", len(errs)))
		return errors.New("部分库存更新失败")
	}

	logger.Info("批量更新库存成功", zap.Int("count", len(updates)))
	return nil
}

// updateStock 更新单个商品库存（带悲观锁）
func (s *ProductService) updateStock(productID uint, quantity int) error {
	return database.Transaction(func(tx *gorm.DB) error {
		var product models.Product

		// 使用FOR UPDATE悲观锁防止并发问题
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, productID).Error; err != nil {
			return err
		}

		// 检查库存是否充足
		newStock := product.Stock + quantity
		if newStock < 0 {
			return errors.New("库存不足")
		}

		// 更新库存
		if err := tx.Model(&product).Update("stock", newStock).Error; err != nil {
			return err
		}

		// 如果库存为0，更新状态
		if newStock == 0 {
			tx.Model(&product).Update("status", "out_of_stock")
		} else if product.Status == "out_of_stock" {
			tx.Model(&product).Update("status", "active")
		}

		return nil
	})
}

// BatchCreateProducts 批量创建商品（Go协程优化）
func (s *ProductService) BatchCreateProducts(products []models.Product) error {
	if len(products) == 0 {
		return errors.New("商品列表为空")
	}

	// 使用协程池批量插入
	const batchSize = 100
	const workerCount = 5

	batches := make(chan []models.Product, (len(products)/batchSize)+1)
	results := make(chan error, (len(products)/batchSize)+1)

	// 启动worker协程池
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range batches {
				// 批量插入
				err := database.DB.CreateInBatches(batch, batchSize).Error
				results <- err
			}
		}()
	}

	// 分批发送任务
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}
		batches <- products[i:end]
	}
	close(batches)

	// 等待所有worker完成
	wg.Wait()
	close(results)

	// 检查是否有错误
	for err := range results {
		if err != nil {
			logger.Error("批量创建商品失败", zap.Error(err))
			return err
		}
	}

	logger.Info("批量创建商品成功", zap.Int("count", len(products)))
	return nil
}

// SearchProducts 全文搜索商品（支持高并发）
func (s *ProductService) SearchProducts(keyword string, page, pageSize int) ([]models.Product, int64, error) {
	if keyword == "" {
		return nil, 0, errors.New("搜索关键词不能为空")
	}

	query := database.DB.Model(&models.Product{}).
		Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Where("status = ?", "active")

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var products []models.Product
	offset := (page - 1) * pageSize
	if err := query.Preload("Category").Offset(offset).Limit(pageSize).Order("sale_count DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
