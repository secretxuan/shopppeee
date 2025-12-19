package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shoppee/ecommerce/internal/config"
	"github.com/shoppee/ecommerce/internal/handler"
	"github.com/shoppee/ecommerce/internal/middleware"
	"github.com/shoppee/ecommerce/internal/websocket"
	"github.com/shoppee/ecommerce/pkg/response"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 设置Gin模式
	if !config.AppConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"app":    config.AppConfig.AppName,
		})
	})

	// 初始化WebSocket
	websocket.InitWebSocket()

	// WebSocket连接（需要认证）
	r.GET("/ws", middleware.AuthMiddleware(), websocket.HandleWebSocket)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 限流中间件（每分钟100次请求）
		// api.Use(middleware.RateLimitMiddleware(100, 1*time.Minute))

		// 认证相关路由（公开）
		authHandler := handler.NewAuthHandler()
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetUserInfo)
		}

		// 商品相关路由（部分公开）
		productHandler := handler.NewProductHandler()
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetProductList)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/:id", productHandler.GetProductByID)

			// 需要管理员权限
			admin := products.Group("")
			admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
			{
				admin.POST("/batch-stock", productHandler.BatchUpdateStock)
			}
		}

		// 用户相关路由（需要认证）
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			// 这里可以添加用户个人中心相关路由
		}

		// 订单相关路由（需要认证）
		order := api.Group("/orders")
		order.Use(middleware.AuthMiddleware())
		{
			// 这里可以添加订单相关路由
		}

		// 购物车相关路由（需要认证）
		cart := api.Group("/cart")
		cart.Use(middleware.AuthMiddleware())
		{
			// 这里可以添加购物车相关路由
		}
	}

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "接口不存在")
	})

	return r
}
