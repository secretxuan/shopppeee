package service

import (
	"fmt"
	"sync"
	"testing"

	"github.com/shoppee/ecommerce/internal/database"
	"github.com/shoppee/ecommerce/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestBatchUpdateStock 测试批量更新库存
func TestBatchUpdateStock(t *testing.T) {
	setupTest()

	productService := NewProductService()

	// 创建测试商品
	products := []models.Product{
		{Name: "商品1", Price: 100, Stock: 100, SKU: "SKU001", Status: "active"},
		{Name: "商品2", Price: 200, Stock: 50, SKU: "SKU002", Status: "active"},
		{Name: "商品3", Price: 300, Stock: 30, SKU: "SKU003", Status: "active"},
	}
	productService.BatchCreateProducts(products)

	tests := []struct {
		name    string
		updates map[uint]int
		wantErr bool
	}{
		{
			name: "正常更新",
			updates: map[uint]int{
				1: -10,
				2: 20,
				3: -5,
			},
			wantErr: false,
		},
		{
			name: "库存不足",
			updates: map[uint]int{
				1: -200, // 超过库存
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := productService.BatchUpdateStock(tt.updates)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBatchUpdateStockConcurrency 测试批量更新库存的并发安全性
func TestBatchUpdateStockConcurrency(t *testing.T) {
	setupTest()

	productService := NewProductService()

	// 创建测试商品
	product := models.Product{
		Name:   "并发测试商品",
		Price:  100,
		Stock:  1000,
		SKU:    "CONCURRENT001",
		Status: "active",
	}
	database.DB.Create(&product)

	// 并发执行100次减库存操作，每次减1
	var wg sync.WaitGroup
	concurrency := 100
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			updates := map[uint]int{product.ID: -1}
			productService.BatchUpdateStock(updates)
		}()
	}

	wg.Wait()

	// 验证最终库存
	var updatedProduct models.Product
	database.DB.First(&updatedProduct, product.ID)
	
	// 库存应该是 1000 - 100 = 900
	assert.Equal(t, 900, updatedProduct.Stock)
}

// BenchmarkBatchUpdateStock 批量更新库存性能基准测试
func BenchmarkBatchUpdateStock(b *testing.B) {
	setupTest()

	productService := NewProductService()

	// 创建100个测试商品
	products := make([]models.Product, 100)
	for i := 0; i < 100; i++ {
		products[i] = models.Product{
			Name:   fmt.Sprintf("商品%d", i),
			Price:  float64(100 + i),
			Stock:  1000,
			SKU:    fmt.Sprintf("SKU%03d", i),
			Status: "active",
		}
	}
	productService.BatchCreateProducts(products)

	// 准备更新数据
	updates := make(map[uint]int)
	for i := uint(1); i <= 100; i++ {
		updates[i] = -1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		productService.BatchUpdateStock(updates)
	}
}
