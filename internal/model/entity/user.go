package entity

import "time"

// 表结构设计

// email / username 唯一索引
// 不存储明文密码

// CREATE TABLE users (
//    id BIGINT PRIMARY KEY AUTO_INCREMENT,
//    email VARCHAR(255) NOT NULL UNIQUE,
//    username VARCHAR(50) NOT NULL UNIQUE,
//    password VARCHAR(255) NOT NULL,
//    bio TEXT,
//    image VARCHAR(255),
//    created_at DATETIME NOT NULL,
//    updated_at DATETIME NOT NULL
//);

type User struct {
	ID       int64  `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex;size:255;not null"`
	Username string `gorm:"uniqueIndex;size:50;not null"`
	Password string `gorm:"size:255;not null"` // 注意这里是hash后的

	Bio       string `gorm:"type:text"`
	Image     string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
