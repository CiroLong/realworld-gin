package router

import (
	"github/CiroLong/realworld-gin/internal/api"
	"github/CiroLong/realworld-gin/internal/middleware"
	"github/CiroLong/realworld-gin/internal/pkg/jwt"
	"github/CiroLong/realworld-gin/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	userService service.UserService,
	articleService service.ArticleService,
	commentService service.CommentService,
	jwtMgr jwt.Manager,
) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// handler
	userHandler := api.NewUserHandler(userService)
	profileHandler := api.NewProfileHandler(userService)
	articleHandler := api.NewArticleHandler(articleService)
	commentHandler := api.NewCommentHandler(commentService)

	// middleware
	auth := middleware.AuthMiddleware(jwtMgr)

	apiGroup := r.Group("/api")

	// ============ Authentication ============
	{
		// POST /api/users - Registration
		apiGroup.POST("/users", userHandler.Register)
		// POST /api/users/login - Authentication
		apiGroup.POST("/users/login", userHandler.Login)
		// GET /api/user - Get current user
		apiGroup.GET("/user", auth, userHandler.GetCurrentUser)
		// PUT /api/user - Update current user
		apiGroup.PUT("/user", auth, userHandler.UpdateCurrentUser)
	}

	// ============ Profiles ============
	{
		// GET /api/profiles/:username - Get profile
		apiGroup.GET("/profiles/:username", profileHandler.GetProfile)
		// POST /api/profiles/:username/follow - Follow user
		apiGroup.POST("/profiles/:username/follow", auth, profileHandler.Follow)
		// DELETE /api/profiles/:username/follow - Unfollow user
		apiGroup.DELETE("/profiles/:username/follow", auth, profileHandler.Unfollow)
	}

	// ============ Articles ============
	{
		// GET /api/articles - List articles
		apiGroup.GET("/articles", articleHandler.ListArticles)
		// GET /api/articles/feed - Feed articles
		apiGroup.GET("/articles/feed", auth, articleHandler.FeedArticles)
		// GET /api/articles/:slug - Get article
		apiGroup.GET("/articles/:slug", articleHandler.GetArticle)
		// POST /api/articles - Create article
		apiGroup.POST("/articles", auth, articleHandler.CreateArticle)
		// PUT /api/articles/:slug - Update article
		apiGroup.PUT("/articles/:slug", auth, articleHandler.UpdateArticle)
		// DELETE /api/articles/:slug - Delete article
		apiGroup.DELETE("/articles/:slug", auth, articleHandler.DeleteArticle)
	}

	// ============ Comments ============
	{
		// POST /api/articles/:slug/comments - Create comment
		apiGroup.POST("/articles/:slug/comments", auth, commentHandler.CreateComment)
		// GET /api/articles/:slug/comments - Get comments
		apiGroup.GET("/articles/:slug/comments", commentHandler.GetComments)
		// DELETE /api/articles/:slug/comments/:id - Delete comment
		apiGroup.DELETE("/articles/:slug/comments/:id", auth, commentHandler.DeleteComment)
	}

	// ============ Favorite ============
	{
		// POST /api/articles/:slug/favorite - Favorite article
		apiGroup.POST("/articles/:slug/favorite", auth, articleHandler.FavoriteArticle)
		// DELETE /api/articles/:slug/favorite - Unfavorite article
		apiGroup.DELETE("/articles/:slug/favorite", auth, articleHandler.UnfavoriteArticle)
	}

	// ============ Tags ============
	{
		// GET /api/tags - Get tags
		apiGroup.GET("/tags", articleHandler.GetTags)
	}

	return r
}
