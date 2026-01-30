package main

import (
	"github/CiroLong/realworld-gin/internal-v2/config"
	"github/CiroLong/realworld-gin/internal-v2/pkg/jwt"
	"github/CiroLong/realworld-gin/internal-v2/repository/gorm"
	"github/CiroLong/realworld-gin/internal-v2/router"
	"github/CiroLong/realworld-gin/internal-v2/service"
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
