package api

import (
	"github.com/gin-gonic/gin"
	"github/CiroLong/realworld-gin/internal/middleware"
	"github/CiroLong/realworld-gin/internal/service"
	"net/http"
)

type ProfileHandler struct {
	userService service.UserService
}

func NewProfileHandler(userService service.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

// GetProfile
// GET /api/profiles/:username
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	username := c.Param("username")

	// userID 是可选的（未登录也能访问）
	var userID int64
	if uid, exists := c.Get(middleware.ContextUserIDKey); exists {
		userID = uid.(int64)
	}

	profile, err := h.userService.GetProfile(
		c.Request.Context(),
		username,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, errError(err))
		return
	}

	c.JSON(http.StatusOK, profile)
}

// Follow
// Authentication required
// POST /api/profiles/:username/follow
func (h *ProfileHandler) Follow(c *gin.Context) {
	username := c.Param("username")

	userIDAny, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, errString("unauthorized"))
		return
	}
	userID := userIDAny.(int64)

	profile, err := h.userService.FollowUserByName(
		c.Request.Context(),
		userID,
		username,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, errError(err))
		return
	}

	c.JSON(http.StatusOK, profile)
}

// Unfollow
// Authentication required
// DELETE /api/profiles/:username/follow
func (h *ProfileHandler) Unfollow(c *gin.Context) {
	username := c.Param("username")

	userIDAny, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, errString("unauthorized"))
		return
	}
	userID := userIDAny.(int64)

	profile, err := h.userService.UnfollowUserByName(
		c.Request.Context(),
		userID,
		username,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, errError(err))
		return
	}

	c.JSON(http.StatusOK, profile)
}
