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
			// 公开接口
			products.GET("", productHandler.GetProductList)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/:id", productHandler.GetProductByID)

			// 需要管理员权限
			admin := products.Group("")
			admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
			{
				admin.POST("", productHandler.CreateProduct)
				admin.PUT("/:id", productHandler.UpdateProduct)
				admin.DELETE("/:id", productHandler.DeleteProduct)
				admin.PATCH("/:id/status", productHandler.UpdateProductStatus)
				admin.POST("/batch-stock", productHandler.BatchUpdateStock)
			}
		}

		// 分类相关路由（部分公开）
		categoryHandler := handler.NewCategoryHandler()
		categories := api.Group("/categories")
		{
			// 公开接口
			categories.GET("", categoryHandler.GetCategoryList)
			categories.GET("/:id", categoryHandler.GetCategory)

			// 需要管理员权限
			admin := categories.Group("")
			admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
			{
				admin.POST("", categoryHandler.CreateCategory)
				admin.PUT("/:id", categoryHandler.UpdateCategory)
				admin.DELETE("/:id", categoryHandler.DeleteCategory)
			}
		}

		// 购物车相关路由（需要认证）
		cartHandler := handler.NewCartHandler()
		cart := api.Group("/cart")
		cart.Use(middleware.AuthMiddleware())
		{
			cart.GET("", cartHandler.GetCart)
			cart.POST("/items", cartHandler.AddToCart)
			cart.PUT("/items/:id", cartHandler.UpdateCartItem)
			cart.DELETE("/items/:id", cartHandler.RemoveFromCart)
			cart.DELETE("/clear", cartHandler.ClearCart)
			cart.PATCH("/items/:id/select", cartHandler.SelectCartItem)
		}

		// 收货地址相关路由（需要认证）
		addressHandler := handler.NewAddressHandler()
		addresses := api.Group("/addresses")
		addresses.Use(middleware.AuthMiddleware())
		{
			addresses.GET("", addressHandler.GetAddressList)
			addresses.GET("/default", addressHandler.GetDefaultAddress)
			addresses.GET("/:id", addressHandler.GetAddress)
			addresses.POST("", addressHandler.CreateAddress)
			addresses.PUT("/:id", addressHandler.UpdateAddress)
			addresses.DELETE("/:id", addressHandler.DeleteAddress)
			addresses.PATCH("/:id/default", addressHandler.SetDefaultAddress)
		}

		// 订单相关路由（需要认证）
		orderHandler := handler.NewOrderHandler()
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetOrderList)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.POST("/:id/cancel", orderHandler.CancelOrder)
			orders.POST("/:id/confirm", orderHandler.ConfirmReceipt)

			// 管理员接口
			admin := orders.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("", orderHandler.AdminGetOrderList)
				admin.PATCH("/:id/status", orderHandler.AdminUpdateOrderStatus)
			}
		}

		// 支付相关路由（需要认证）
		paymentHandler := handler.NewPaymentHandler()
		payments := api.Group("/payments")
		payments.Use(middleware.AuthMiddleware())
		{
			payments.POST("", paymentHandler.CreatePayment)
			payments.GET("/:id", paymentHandler.GetPayment)
		}
		// 支付回调（公开接口）
		api.POST("/payment-callback", paymentHandler.PaymentCallback)

		// 评价相关路由
		reviewHandler := handler.NewReviewHandler()
		reviews := api.Group("/reviews")
		{
			// 公开接口
			reviews.GET("/products/:id", reviewHandler.GetProductReviews)

			// 需要认证
			auth := reviews.Group("")
			auth.Use(middleware.AuthMiddleware())
			{
				auth.POST("", reviewHandler.CreateReview)
				auth.GET("/my", reviewHandler.GetMyReviews)
				auth.DELETE("/:id", reviewHandler.DeleteReview)

				// 管理员接口
				admin := auth.Group("")
				admin.Use(middleware.AdminMiddleware())
				{
					admin.POST("/:id/reply", reviewHandler.AdminReplyReview)
				}
			}
		}
	}

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "接口不存在")
	})

	return r
}
