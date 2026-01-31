package entity

// CREATE TABLE favorites (
//  user_id BIGINT NOT NULL,
//  article_id BIGINT NOT NULL,
//
//  PRIMARY KEY (user_id, article_id),
//
//  CONSTRAINT fk_favorites_user
//    FOREIGN KEY (user_id) REFERENCES users(id),
//  CONSTRAINT fk_favorites_article
//    FOREIGN KEY (article_id) REFERENCES articles(id)
//);

// 注意 这里不在 struct 里显式声明 GORM 的外键关系

type Favorite struct {
	UserID    int64 `gorm:"primaryKey"`
	ArticleID int64 `gorm:"primaryKey"`
}
