package gorm

import (
	"context"
	"github/CiroLong/realworld-gin/internal/model/entity"
	"github/CiroLong/realworld-gin/internal/repository"

	"gorm.io/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) repository.CommentRepo {
	return &CommentRepo{db: db}
}

func (c CommentRepo) Create(ctx context.Context, comment *entity.Comment) error {
	// TODO: check id time是否正确填写
	return c.db.WithContext(ctx).Create(comment).Error
}

func (c CommentRepo) ListByArticle(ctx context.Context, articleID int64) ([]*entity.Comment, error) {
	var comments []*entity.Comment

	err := c.db.WithContext(ctx).
		Where("article_id = ?", articleID).
		Order("created_at ASC").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c CommentRepo) FindByID(ctx context.Context, id int64) (*entity.Comment, error) {
	var comment entity.Comment

	err := c.db.WithContext(ctx).
		First(&comment, id).Error

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (c CommentRepo) Delete(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).
		Delete(&entity.Comment{}, id).Error
}
