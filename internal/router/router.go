package router

import (
	"github.com/gin-gonic/gin"
	"github/CiroLong/realworld-gin/internal/api"
	"github/CiroLong/realworld-gin/internal/middleware"
	"github/CiroLong/realworld-gin/internal/pkg/jwt"
	"github/CiroLong/realworld-gin/internal/service"
)

func NewRouter(userService service.UserService, jwtMgr jwt.Manager) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// handler
	userHandler := api.NewUserHandler(userService)

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

	return r
}
