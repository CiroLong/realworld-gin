package router

import (
	"github.com/gin-gonic/gin"
	"github/CiroLong/realworld-gin/internal/api"
	"github/CiroLong/realworld-gin/internal/middleware"
	"github/CiroLong/realworld-gin/internal/pkg/jwt"
	"github/CiroLong/realworld-gin/internal/service"
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
	articleHandler := api.NewArticleHandler(articleService)
	commentHandler := api.NewCommentHandler(commentService)

	apiGroup := r.Group("/api")
	{
		// 注册 & 登录（公开）
		apiGroup.POST("/users", userHandler.Register)
		apiGroup.POST("/users/login", userHandler.Login)

		// 需要登录
		authGroup := apiGroup.Group("/")
		authGroup.Use(middleware.AuthMiddleware(jwtMgr))
		{
			authGroup.GET("/user", userHandler.GetCurrentUser)
			authGroup.PUT("/user", userHandler.UpdateCurrentUser)
		}
	}
	{
		articles := apiGroup.Group("/articles")
		// -------- tags --------
		articles.GET("/tags", articleHandler.GetTags)
		// -------- feed --------
		articles.GET("/feed", middleware.AuthMiddleware(jwtMgr), articleHandler.FeedArticles)
		// -------- public --------
		articles.GET("", articleHandler.ListArticles)
		articles.GET("/:slug", articleHandler.GetArticle)

		// -------- auth required --------
		articlesAuthGroup := articles.Use(middleware.AuthMiddleware(jwtMgr))
		{
			articlesAuthGroup.POST("", articleHandler.CreateArticle)
			articlesAuthGroup.PUT("/:slug", articleHandler.UpdateArticle)
			articlesAuthGroup.DELETE("/:slug", articleHandler.DeleteArticle)

			articlesAuthGroup.POST("/:slug/favorite", articleHandler.FavoriteArticle)
			articlesAuthGroup.DELETE("/:slug/favorite", articleHandler.UnfavoriteArticle)
		}
	}
	comments := r.Group("/api/articles/:slug/comments")

	{
		comments.GET("", commentHandler.GetComments)
		comments.POST("", middleware.AuthMiddleware(jwtMgr), commentHandler.CreateComment)
		comments.DELETE("/:id", middleware.AuthMiddleware(jwtMgr), commentHandler.DeleteComment)
	}

	return r
}
