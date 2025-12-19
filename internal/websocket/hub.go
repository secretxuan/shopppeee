package websocket

import (
	"encoding/json"
	"sync"

	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

// Hub WebSocket连接中心（管理所有客户端连接）
type Hub struct {
	// 注册的客户端连接
	clients map[*Client]bool

	// 用户ID到客户端的映射（支持按用户推送）
	userClients map[uint]*Client

	// 广播消息通道
	broadcast chan []byte

	// 注册请求通道
	register chan *Client

	// 注销请求通道
	unregister chan *Client

	// 互斥锁
	mu sync.RWMutex
}

// Message WebSocket消息结构
type Message struct {
	Type    string      `json:"type"`    // message_type: system, order, notification
	Content interface{} `json:"content"` // 消息内容
	UserID  uint        `json:"user_id,omitempty"`
	Time    int64       `json:"time"`
}

// NewHub 创建新的Hub实例
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		userClients: make(map[uint]*Client),
		broadcast:   make(chan []byte, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
}

// Run 运行Hub（持续监听各个通道）
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if client.userID > 0 {
				h.userClients[client.userID] = client
			}
			h.mu.Unlock()
			logger.Info("WebSocket客户端已连接",
				zap.Uint("user_id", client.userID),
				zap.Int("total_clients", len(h.clients)),
			)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.userID > 0 {
					delete(h.userClients, client.userID)
				}
				close(client.send)
			}
			h.mu.Unlock()
			logger.Info("WebSocket客户端已断开",
				zap.Uint("user_id", client.userID),
				zap.Int("total_clients", len(h.clients)),
			)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// 发送失败，关闭客户端连接
					close(client.send)
					delete(h.clients, client)
					if client.userID > 0 {
						delete(h.userClients, client.userID)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastMessage 广播消息给所有客户端
func (h *Hub) BroadcastMessage(msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("序列化消息失败", zap.Error(err))
		return
	}
	h.broadcast <- data
	logger.Debug("广播消息", zap.String("type", msg.Type))
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID uint, msg *Message) bool {
	h.mu.RLock()
	client, exists := h.userClients[userID]
	h.mu.RUnlock()

	if !exists {
		logger.Warn("用户未连接", zap.Uint("user_id", userID))
		return false
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("序列化消息失败", zap.Error(err))
		return false
	}

	select {
	case client.send <- data:
		logger.Debug("发送消息给用户", zap.Uint("user_id", userID), zap.String("type", msg.Type))
		return true
	default:
		logger.Warn("发送消息失败，通道已满", zap.Uint("user_id", userID))
		return false
	}
}

// GetOnlineUserCount 获取在线用户数
func (h *Hub) GetOnlineUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.userClients)
}

// GetTotalConnections 获取总连接数
func (h *Hub) GetTotalConnections() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
