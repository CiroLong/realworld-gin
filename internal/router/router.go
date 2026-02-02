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

	// ==================== Authentication ====================
	// 认证相关路由
	usersGroup := apiGroup.Group("/users")
	{
		// 公开路由
		usersGroup.POST("", userHandler.Register)         // POST /api/users - 注册
		usersGroup.POST("/login", userHandler.Login)    // POST /api/users/login - 登录
	}

	// 当前用户相关路由（需要认证）
	userGroup := apiGroup.Group("/user")
	userGroup.Use(auth)
	{
		userGroup.GET("", userHandler.GetCurrentUser)    // GET /api/user - 获取当前用户
		userGroup.PUT("", userHandler.UpdateCurrentUser) // PUT /api/user - 更新当前用户
	}

	// ==================== Profiles ====================
	// 用户资料相关路由
	profilesGroup := apiGroup.Group("/profiles")
	{
		// 公开路由
		profilesGroup.GET("/:username", profileHandler.GetProfile) // GET /api/profiles/:username - 获取用户资料

		// 需要认证的路由
		profilesAuthGroup := profilesGroup.Group("/:username")
		profilesAuthGroup.Use(auth)
		{
			profilesAuthGroup.POST("/follow", profileHandler.Follow)       // POST /api/profiles/:username/follow - 关注用户
			profilesAuthGroup.DELETE("/follow", profileHandler.Unfollow)    // DELETE /api/profiles/:username/follow - 取消关注
		}
	}

	// ==================== Articles ====================
	// 文章相关路由
	articlesGroup := apiGroup.Group("/articles")
	{
		// 公开路由
		articlesGroup.GET("", articleHandler.ListArticles)         // GET /api/articles - 文章列表

		// Feed 路由（需要认证）- 必须放在 /:slug 前面，否则会被当作 slug 处理
		articlesGroup.GET("/feed", auth, articleHandler.FeedArticles)   // GET /api/articles/feed - 文章Feed

		// 公开路由
		articlesGroup.GET("/:slug", articleHandler.GetArticle)     // GET /api/articles/:slug - 获取文章详情

		// 需要认证的路由
		articlesAuthGroup := articlesGroup.Group("")
		articlesAuthGroup.Use(auth)
		{
			articlesAuthGroup.POST("", articleHandler.CreateArticle)                 // POST /api/articles - 创建文章
			articlesAuthGroup.PUT("/:slug", articleHandler.UpdateArticle)            // PUT /api/articles/:slug - 更新文章
			articlesAuthGroup.DELETE("/:slug", articleHandler.DeleteArticle)         // DELETE /api/articles/:slug - 删除文章
			articlesAuthGroup.POST("/:slug/favorite", articleHandler.FavoriteArticle)     // POST /api/articles/:slug/favorite - 收藏文章
			articlesAuthGroup.DELETE("/:slug/favorite", articleHandler.UnfavoriteArticle) // DELETE /api/articles/:slug/favorite - 取消收藏
		}
	}

	// ==================== Comments ====================
	// 评论相关路由
	commentsGroup := apiGroup.Group("/articles/:slug/comments")
	{
		// 公开路由
		commentsGroup.GET("", commentHandler.GetComments) // GET /api/articles/:slug/comments - 获取评论

		// 需要认证的路由
		commentsAuthGroup := commentsGroup.Group("")
		commentsAuthGroup.Use(auth)
		{
			commentsAuthGroup.POST("", commentHandler.CreateComment)           // POST /api/articles/:slug/comments - 创建评论
			commentsAuthGroup.DELETE("/:id", commentHandler.DeleteComment)     // DELETE /api/articles/:slug/comments/:id - 删除评论
		}
	}

	// ==================== Tags ====================
	// 标签相关路由（公开）
	apiGroup.GET("/tags", articleHandler.GetTags) // GET /api/tags - 获取标签列表

	return r
}
