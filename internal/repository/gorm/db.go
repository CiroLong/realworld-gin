package gorm

import (
	"github/CiroLong/realworld-gin/internal/config"
	"github/CiroLong/realworld-gin/internal/model/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := config.C().Database

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	var err error
	DB, err = gorm.Open(mysql.Open(cfg.DSN), gormCfg)
	if err != nil {
		return err
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 连接池配置（非常重要）
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// AutoMigrate
func AutoMigrate() error {
	return DB.AutoMigrate(
		&entity.User{},
		// TODO: 加入剩余entity
		// &entity.Article{},
		// &entity.Comment{},
	)
}
