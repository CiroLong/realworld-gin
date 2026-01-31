package entity

// CREATE TABLE tags (
//  id BIGINT AUTO_INCREMENT PRIMARY KEY,
//  name VARCHAR(50) NOT NULL UNIQUE
//);

type Tag struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"size:50;uniqueIndex;not null"`
}
