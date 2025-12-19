package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shoppee/ecommerce/internal/database"
	"github.com/shoppee/ecommerce/internal/models"
	"github.com/shoppee/ecommerce/pkg/jwt"
	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct{}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Phone    string `json:"phone"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expires_at"`
	User      *models.User `json:"user"`
}

// Register 用户注册（支持高并发）
func (s *AuthService) Register(req *RegisterRequest) (*models.User, error) {
	// 检查用户名是否存在
	var existUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否存在
	if err := database.DB.Where("email = ?", req.Email).First(&existUser).Error; err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // BeforeCreate钩子会自动加密
		Phone:    req.Phone,
		Role:     "user",
		Status:   "active",
	}

	// 使用事务创建用户和购物车
	err := database.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 创建购物车
		cart := &models.Cart{UserID: user.ID}
		if err := tx.Create(cart).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error("用户注册失败", zap.Error(err))
		return nil, err
	}

	logger.Info("用户注册成功", zap.Uint("user_id", user.ID))
	return user, nil
}

// Login 用户登录（高并发优化：Redis缓存）
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查询用户
	var user models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查账户状态
	if user.Status != "active" {
		return nil, errors.New("账户已被禁用")
	}

	// 生成token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		logger.Error("生成token失败", zap.Error(err))
		return nil, err
	}

	// 更新最后登录时间（异步执行，不阻塞响应）
	go func() {
		now := time.Now()
		database.DB.Model(&user).Update("last_login", now)

		// 缓存用户信息到Redis（7天过期）
		ctx := context.Background()
		userKey := fmt.Sprintf("user:%d", user.ID)
		database.RedisClient.Set(ctx, userKey, user.Username, 7*24*time.Hour)
	}()

	// 隐藏敏感信息
	user.Password = ""

	logger.Info("用户登录成功", zap.Uint("user_id", user.ID))

	return &LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		User:      &user,
	}, nil
}

// GetUserInfo 获取用户信息（带Redis缓存）
func (s *AuthService) GetUserInfo(userID uint) (*models.User, error) {
	var user models.User

	// 先从Redis获取
	ctx := context.Background()
	userKey := fmt.Sprintf("user:%d", userID)
	
	_, err := database.RedisClient.Get(ctx, userKey).Result()
	if err == nil {
		// 缓存命中，从数据库加载完整信息
		if err := database.DB.Preload("Addresses").First(&user, userID).Error; err != nil {
			return nil, err
		}
		user.Password = ""
		return &user, nil
	}

	// 缓存未命中，从数据库查询
	if err := database.DB.Preload("Addresses").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 异步缓存到Redis
	go func() {
		database.RedisClient.Set(ctx, userKey, user.Username, 7*24*time.Hour)
	}()

	user.Password = ""
	return &user, nil
}
