package websocket

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shoppee/ecommerce/internal/middleware"
	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

var (
	// GlobalHub 全局WebSocket Hub实例
	GlobalHub *Hub

	// upgrader WebSocket升级器
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有跨域请求（生产环境应配置白名单）
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// InitWebSocket 初始化WebSocket服务
func InitWebSocket() {
	GlobalHub = NewHub()
	go GlobalHub.Run()
	logger.Info("WebSocket服务已启动")
}

// HandleWebSocket 处理WebSocket连接
func HandleWebSocket(c *gin.Context) {
	// 获取当前用户ID（从JWT认证中间件）
	userID := middleware.GetCurrentUserID(c)

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("WebSocket升级失败", zap.Error(err))
		return
	}

	// 创建客户端
	client := &Client{
		hub:    GlobalHub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}

	// 注册客户端
	client.hub.register <- client

	// 发送欢迎消息
	welcomeMsg := &Message{
		Type:    "system",
		Content: "欢迎连接Shoppee实时消息服务",
		Time:    time.Now().Unix(),
	}
	GlobalHub.SendToUser(userID, welcomeMsg)

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}

// NotifyOrderStatus 通知订单状态变更
func NotifyOrderStatus(userID uint, orderID uint, status string) {
	if GlobalHub == nil {
		return
	}

	msg := &Message{
		Type: "order",
		Content: map[string]interface{}{
			"order_id": orderID,
			"status":   status,
			"message":  "您的订单状态已更新",
		},
		UserID: userID,
		Time:   time.Now().Unix(),
	}

	GlobalHub.SendToUser(userID, msg)
}

// BroadcastPromotion 广播促销信息
func BroadcastPromotion(title, content string) {
	if GlobalHub == nil {
		return
	}

	msg := &Message{
		Type: "promotion",
		Content: map[string]interface{}{
			"title":   title,
			"content": content,
		},
		Time: time.Now().Unix(),
	}

	GlobalHub.BroadcastMessage(msg)
}

// NotifyStockAlert 库存预警通知
func NotifyStockAlert(adminUserID uint, productID uint, productName string, stock int) {
	if GlobalHub == nil {
		return
	}

	msg := &Message{
		Type: "stock_alert",
		Content: map[string]interface{}{
			"product_id":   productID,
			"product_name": productName,
			"stock":        stock,
			"message":      "商品库存不足，请及时补货",
		},
		UserID: adminUserID,
		Time:   time.Now().Unix(),
	}

	GlobalHub.SendToUser(adminUserID, msg)
}
