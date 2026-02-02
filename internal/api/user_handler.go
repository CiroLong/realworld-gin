package api

import (
	"github/CiroLong/realworld-gin/internal/middleware"
	"github/CiroLong/realworld-gin/internal/model/dto"
	"github/CiroLong/realworld-gin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// 下面是user相关接口实现

// Register
// Registration
// POST /api/users
func (h *UserHandler) Register(c *gin.Context) {
	// 1. 处理请求
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": gin.H{"body": []string{err.Error()}},
		})
		return
	}

	// 2. 调service注册
	resp, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": gin.H{"body": []string{err.Error()}},
		})
		return
	}

	// 3. 返回
	c.JSON(http.StatusOK, resp)
}

// Login
// Authentication: 登陆/认证
// POST /api/users/login
func (h *UserHandler) Login(c *gin.Context) {
	// 1. 处理请求
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": gin.H{"body": []string{err.Error()}},
		})
		return
	}

	// 2. 调service认证
	resp, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": gin.H{"body": []string{err.Error()}},
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCurrentUser
// Auth needed
// GET /api/user
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// 1.从context中获取userID
	uidVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": gin.H{"body": []string{"unauthorized"}},
		})
		return
	}

	userID := uidVal.(int64)

	// 2. 调service处理逻辑
	resp, err := h.userService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": gin.H{"body": []string{err.Error()}},
		})
		return
	}

	// 3. 返回
	c.JSON(http.StatusOK, resp)
}

// UpdateCurrentUser
// PUT /api/user
// auth needed
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	// 1. 从context中取userID
	uidVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": gin.H{"body": []string{"unauthorized"}},
		})
		return
	}
	userID := uidVal.(int64)

	// 2. 解析请求
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": gin.H{"body": []string{err.Error()}},
		})
		return
	}

	// 3. 处理
	resp, err := h.userService.UpdateCurrentUser(
		c.Request.Context(),
		userID,
		&req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": gin.H{"body": []string{err.Error()}},
		})
		return
	}

	// 4. 返回
	c.JSON(http.StatusOK, resp)
}
