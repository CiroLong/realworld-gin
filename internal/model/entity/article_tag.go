package entity

// 关系表(多对多)

// CREATE TABLE article_tags (
//  article_id BIGINT NOT NULL,
//  tag_id BIGINT NOT NULL,
//
//  PRIMARY KEY (article_id, tag_id),
//
//  CONSTRAINT fk_article_tags_article
//    FOREIGN KEY (article_id) REFERENCES articles(id),
//  CONSTRAINT fk_article_tags_tag
//    FOREIGN KEY (tag_id) REFERENCES tags(id)
//);

// 注意 这里不在 struct 里显式声明 GORM 的外键关系

type ArticleTag struct {
	ArticleID int64 `gorm:"primaryKey"`
	TagID     int64 `gorm:"primaryKey"`
}
