package service

import (
	"fmt"
	"testing"

	"github.com/shoppee/ecommerce/internal/config"
	"github.com/shoppee/ecommerce/internal/database"
	"github.com/stretchr/testify/assert"
)

// 测试前的初始化
func setupTest() {
	// 初始化测试配置
	config.InitConfig()
	
	// 初始化测试数据库（可使用SQLite内存数据库）
	database.InitDB()
	database.AutoMigrate()
}

// TestRegister 测试用户注册
func TestRegister(t *testing.T) {
	setupTest()

	authService := NewAuthService()

	tests := []struct {
		name    string
		req     *RegisterRequest
		wantErr bool
	}{
		{
			name: "正常注册",
			req: &RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Phone:    "13800138000",
			},
			wantErr: false,
		},
		{
			name: "用户名重复",
			req: &RegisterRequest{
				Username: "testuser",
				Email:    "test2@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "邮箱重复",
			req: &RegisterRequest{
				Username: "testuser2",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := authService.Register(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.req.Username, user.Username)
				assert.Equal(t, tt.req.Email, user.Email)
			}
		})
	}
}

// TestLogin 测试用户登录
func TestLogin(t *testing.T) {
	setupTest()

	authService := NewAuthService()

	// 先注册一个用户
	registerReq := &RegisterRequest{
		Username: "logintest",
		Email:    "logintest@example.com",
		Password: "password123",
	}
	authService.Register(registerReq)

	tests := []struct {
		name    string
		req     *LoginRequest
		wantErr bool
	}{
		{
			name: "正确的用户名和密码",
			req: &LoginRequest{
				Username: "logintest",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "错误的密码",
			req: &LoginRequest{
				Username: "logintest",
				Password: "wrongpassword",
			},
			wantErr: true,
		},
		{
			name: "不存在的用户",
			req: &LoginRequest{
				Username: "notexist",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := authService.Login(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Token)
				assert.NotNil(t, resp.User)
			}
		})
	}
}

// BenchmarkRegister 注册性能基准测试
func BenchmarkRegister(b *testing.B) {
	setupTest()
	authService := NewAuthService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &RegisterRequest{
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: "password123",
		}
		authService.Register(req)
	}
}
