package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shoppee/ecommerce/internal/middleware"
	"github.com/shoppee/ecommerce/internal/service"
	"github.com/shoppee/ecommerce/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=models.User}
// @Failure 400 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 隐藏密码
	user.Password = ""
	response.SuccessWithMessage(c, "注册成功", user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取JWT token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=service.LoginResponse}
// @Failure 400 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	loginResp, err := h.authService.Login(&req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, loginResp)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取已登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.User}
// @Failure 401 {object} response.Response
// @Router /auth/me [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "未登录")
		return
	}

	user, err := h.authService.GetUserInfo(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	response.Success(c, user)
}
