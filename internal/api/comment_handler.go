package api

import (
	"github/CiroLong/realworld-gin/internal/middleware"
	"github/CiroLong/realworld-gin/internal/model/dto"
	"github/CiroLong/realworld-gin/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment
// Add Comments to an Article
// Authentication required
// POST /api/articles/:slug/comments
func (h *CommentHandler) CreateComment(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, errString("empty slug"))
		return
	}

	userIDany, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, errString("unauthorized"))
		return
	}
	userID := userIDany.(int64)

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errError(err))
		return
	}

	resp, err := h.commentService.CreateComment(
		c.Request.Context(),
		userID,
		slug,
		&req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, errError(err))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetComments
// GET /api/articles/:slug/comments
func (h *CommentHandler) GetComments(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, errString("empty slug"))
		return
	}

	var userID int64 = 0
	userIDany, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		userID = userIDany.(int64)
	}

	resp, err := h.commentService.GetComments(
		c.Request.Context(),
		slug,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, errError(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteComment
// DELETE /api/articles/:slug/comments/:id
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	if commentIDStr == "" {
		c.JSON(http.StatusBadRequest, errString("empty comment id"))
		return
	}

	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errString("invalid comment id"))
		return
	}

	userIDany, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, errString("unauthorized"))
		return
	}
	userID := userIDany.(int64)

	if err := h.commentService.DeleteComment(
		c.Request.Context(),
		userID,
		commentID,
	); err != nil {
		c.JSON(http.StatusForbidden, errError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
