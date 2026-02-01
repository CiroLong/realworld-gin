package repository

import (
	"context"
	"github/CiroLong/realworld-gin/internal/model/entity"
)

type CommentRepo interface {
	// Create 创建评论
	Create(ctx context.Context, comment *entity.Comment) error

	// ListByArticle 获取文章下所有评论（按创建时间正序）
	ListByArticle(ctx context.Context, articleID int64) ([]*entity.Comment, error)

	// FindByID 根据 comment id 查找
	FindByID(ctx context.Context, id int64) (*entity.Comment, error)

	// Delete 删除评论（物理删除）
	Delete(ctx context.Context, id int64) error
}
