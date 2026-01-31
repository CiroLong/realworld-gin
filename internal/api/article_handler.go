package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github/CiroLong/realworld-gin/internal/middleware"
	"github/CiroLong/realworld-gin/internal/model/dto"
	"github/CiroLong/realworld-gin/internal/repository"
	"github/CiroLong/realworld-gin/internal/service"
	"net/http"
	"strconv"
)

type ArticleHandler struct {
	articleService service.ArticleService
}

func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// 下面是article相关接口实现

// CreateArticle
// Authentication required
// POST /api/articles
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	// 1. 绑定请求 DTO
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errError(err))
		return
	}

	// 2. 获取当前登录用户 ID（AuthMiddleware 已设置）
	userID, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, errString("unauthorized"))
		return
	}

	// 3. 调用 Service
	articleResp, err := h.articleService.CreateArticle(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errError(err))
		return
	}

	c.JSON(http.StatusCreated, articleResp)
}

// GetArticle
// GET /api/articles/:slug
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	// 1. 绑定参数 slug
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, errString("slug cannot be empty"))
		return
	}

	// 2. 调用service
	resp, err := h.articleService.GetArticle(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, errError(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errError(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateArticle
// Authentication required
// PUT /api/articles/:slug
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	// 1.参数绑定
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, errString("slug cannot be empty"))
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errError(err))
		return
	}

	// 2. 取userID
	userID, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, errString("unauthorized"))
		return
	}

	// 3. 调用service
	resp, err := h.articleService.UpdateArticle(c.Request.Context(), slug, userID.(int64), &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, errError(err))
			return
		}
		if errors.Is(err, service.ErrPermissionDenied) {
			c.JSON(http.StatusForbidden, errError(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errError(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteArticle
// Authentication required
// DELETE /api/articles/:slug
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, errString("slug cannot be empty"))
		return
	}

	userID, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, errString("unauthorized"))
		return
	}

	err := h.articleService.DeleteArticle(c.Request.Context(), slug, userID.(int64))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, errError(err))
			return
		}
		if errors.Is(err, service.ErrPermissionDenied) {
			c.JSON(http.StatusForbidden, errError(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// FavoriteArticle
// Authentication required
// POST /api/articles/:slug/favorite
func (h *ArticleHandler) FavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, errString("slug cannot be empty"))
		return
	}

	userID, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, errString("unauthorized"))
		return
	}

	resp, err := h.articleService.FavoriteArticle(c.Request.Context(), slug, userID.(int64))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, errError(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errError(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UnfavoriteArticle
// Authentication required
// DELETE /api/articles/:slug/favorite
func (h *ArticleHandler) UnfavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, errString("slug cannot be empty"))
		return
	}

	userID, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, errString("unauthorized"))
		return
	}

	resp, err := h.articleService.UnfavoriteArticle(c.Request.Context(), slug, userID.(int64))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, errError(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errError(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListArticles
// Authentication optional
// GET /api/articles?tag=&author=&favorited=&limit=&offset=
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	tag := c.Query("tag")
	author := c.Query("author")
	favorited := c.Query("favorited")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var userID int64
	if uid, exists := c.Get(middleware.ContextUserIDKey); exists {
		userID = int64(uid.(uint))
	}

	resp, err := h.articleService.ListArticles(c.Request.Context(), tag, author, favorited, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errError(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// FeedArticles
// Authentication required
// GET /api/articles/feed?limit=&offset=
func (h *ArticleHandler) FeedArticles(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, errString("unauthorized"))
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.articleService.FeedArticles(c.Request.Context(), int64(userID.(uint)), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errError(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTags
// GET /api/tags
func (h *ArticleHandler) GetTags(c *gin.Context) {
	tags, err := h.articleService.ListTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
