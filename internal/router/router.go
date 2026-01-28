package router

import (
	"github.com/gin-gonic/gin"
	"github/CiroLong/realworld-gin/internal/controller"
	"github/CiroLong/realworld-gin/internal/middleware"
)

// 只负责把各个分组（/api/users、/api/articles）同对应的 handler “连线” 并挂载中间件。

func Register(r *gin.Engine) {
	uc := controller.NewUserController()
	r.POST("/users", uc.RegisterUsers)
	r.POST("/users/login", uc.Login)

	// Authentication required
	r.Use(middleware.AuthMiddleware(true))
	r.GET("/user", uc.GetCurrentUser)
}
