package repository

import (
	"github/CiroLong/realworld-gin/internal/common"
	"gorm.io/gorm"
)

type ArticleRepository interface {
}
type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository() ArticleRepository {
	return &articleRepository{db: common.GetDB()}
}
