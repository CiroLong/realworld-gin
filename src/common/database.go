package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

func InitDB() *gorm.DB {

	dsn := "root:123456@tcp(127.0.0.1:3306)/realworld?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println("db err: (Init) ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("sqlDB get error:", err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return DB
}
