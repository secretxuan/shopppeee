package websocket

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

const (
	// 允许向客户端写入消息的时间
	writeWait = 10 * time.Second

	// 允许从客户端读取下一个pong消息的时间
	pongWait = 60 * time.Second

	// 在此期间向客户端发送ping，必须小于pongWait
	pingPeriod = (pongWait * 9) / 10

	// 允许的最大消息大小
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client WebSocket客户端
type Client struct {
	// Hub引用
	hub *Hub

	// WebSocket连接
	conn *websocket.Conn

	// 发送消息的缓冲通道
	send chan []byte

	// 用户ID
	userID uint
}

// readPump 从WebSocket连接读取消息并转发到Hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("WebSocket读取错误", zap.Error(err))
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// 处理客户端发送的消息
		c.handleMessage(message)
	}
}

// writePump 从Hub接收消息并写入WebSocket连接
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub关闭了通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将队列中的消息一起发送
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理客户端发送的消息
func (c *Client) handleMessage(message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		logger.Error("解析消息失败", zap.Error(err))
		return
	}

	// 根据消息类型处理
	switch msg.Type {
	case "ping":
		// 响应心跳
		response := Message{
			Type:    "pong",
			Content: "ok",
			Time:    time.Now().Unix(),
		}
		data, _ := json.Marshal(response)
		c.send <- data

	case "echo":
		// 回显消息
		msg.Time = time.Now().Unix()
		data, _ := json.Marshal(msg)
		c.send <- data

	default:
		logger.Debug("收到客户端消息",
			zap.Uint("user_id", c.userID),
			zap.String("type", msg.Type),
		)
	}
}
