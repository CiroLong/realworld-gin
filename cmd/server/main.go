package main

import (
	"fmt"
	"github/CiroLong/realworld-gin/internal/config"
	"github/CiroLong/realworld-gin/internal/pkg/jwt"
	"github/CiroLong/realworld-gin/internal/repository/gorm"
	"github/CiroLong/realworld-gin/internal/router"
	"github/CiroLong/realworld-gin/internal/service"
	"log"
)

// 	运行流程
//	从 main.go
//	加载配置
//	初始化数据库连接
//	依赖注入
//	注册路由和中间件
//	启动 Gin 服务

func main() {
	// 1. 读配置
	if err := config.Load(); err != nil {
		log.Fatalf("load cfg failed: %v", err)
	}
	cfg := config.C()
	fmt.Printf("%#v", cfg)

	// 2. 链接数据库
	if err := gorm.InitDB(); err != nil {
		log.Fatalf("initDB failed: %v", err)
	}
	db := gorm.GetDB()
	err := gorm.AutoMigrate()
	if err != nil {
		log.Fatalf("initDB AutoMigrate failed: %v", err)
	}

	// 3. 参数注入service
	jwtMgr := jwt.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.ExpireTime,
	)

	userRepo := gorm.NewUserRepo(db)
	userService := service.NewUserService(userRepo, jwtMgr)

	// 4. 注册路由和中间件
	r := router.NewRouter(userService, jwtMgr)
	r.Run(":8080")
}
