package entity

import "time"

type Comment struct {
	ID int64 `gorm:"primaryKey"`

	Body string `gorm:"type:text;not null"`

	ArticleID int64 `gorm:"index;not null"`
	AuthorID  int64 `gorm:"index;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
