package entity

import "time"

// 表结构设计

// CREATE TABLE articles (
//  id BIGINT AUTO_INCREMENT PRIMARY KEY,
//  slug VARCHAR(255) NOT NULL UNIQUE,
//  title VARCHAR(255) NOT NULL,
//  description VARCHAR(255) NOT NULL,
//  body TEXT NOT NULL,
//
//  author_id BIGINT NOT NULL,
//
//  favorites_count INT NOT NULL DEFAULT 0,
//
//  created_at DATETIME NOT NULL,
//  updated_at DATETIME NOT NULL,
//
//  INDEX idx_author_id (author_id),
//  CONSTRAINT fk_articles_author
//    FOREIGN KEY (author_id) REFERENCES users(id)
//);

// 注意 这里不在 struct 里显式声明 GORM 的外键关系

type Article struct {
	ID          int64  `gorm:"primaryKey"`
	Slug        string `gorm:"size:255;uniqueIndex;not null"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"size:255;not null"`
	Body        string `gorm:"type:text;not null"`

	AuthorID int64 `gorm:"index;not null"`

	FavoritesCount int `gorm:"not null;default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
