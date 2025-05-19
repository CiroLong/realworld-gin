package router

import (
	"github.com/gin-gonic/gin"
	"github/CiroLong/realworld-gin/internal/controller"
)

// 只负责把各个分组（/api/users、/api/articles）同对应的 handler “连线” 并挂载中间件。

func Register(r *gin.Engine) {
	r.POST("/users", controller.NewUserController().RegisterUsers)
}
