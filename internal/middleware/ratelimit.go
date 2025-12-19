package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shoppee/ecommerce/internal/database"
	"github.com/shoppee/ecommerce/pkg/response"
)

// RateLimitMiddleware 限流中间件（基于Redis令牌桶算法）
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP作为限流key
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", clientIP)

		ctx := context.Background()

		// 使用Redis实现滑动窗口限流
		now := time.Now().Unix()
		windowStart := now - int64(window.Seconds())

		pipe := database.RedisClient.Pipeline()

		// 移除窗口外的记录
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

		// 添加当前请求
		pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", now)})

		// 统计窗口内请求数
		pipe.ZCard(ctx, key)

		// 设置过期时间
		pipe.Expire(ctx, key, window)

		cmds, err := pipe.Exec(ctx)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "限流检查失败")
			c.Abort()
			return
		}

		// 获取请求数
		count := cmds[2].(*redis.IntCmd).Val()

		if int(count) > limit {
			response.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
