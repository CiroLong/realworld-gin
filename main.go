package main

import (
	"github.com/gin-gonic/gin"
	"github/CiroLong/realworld-gin/internal/common"
	"github/CiroLong/realworld-gin/internal/repository"
)

// 	运行流程
//	从 main.go 加载配置
//	初始化数据库连接
//	通过依赖注入逐层构建 Controller
//	注册路由和中间件
//	启动 Gin 服务

func main() {
	// Do some work to init
	// Load Config

	// Init DB
	common.InitDB()
	repository.AutoMigrate()

	// BuildController

	// register router
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// middleware

	// Run

	router.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
