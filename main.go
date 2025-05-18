package main

import (
	"github.com/gin-gonic/gin"
	"github/CiroLong/realworld-gin/src/common"
	"github/CiroLong/realworld-gin/src/models"
)

func main() {

	// Do some work to init
	common.InitDB()
	models.AutoMigrate()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
