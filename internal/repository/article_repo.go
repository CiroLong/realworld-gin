package repository

import (
	"context"
	"github/CiroLong/realworld-gin/internal/model/entity"
)

type ListArticlesFilter struct {
	Tag         *string
	Author      *string
	FavoritedBy *string

	Limit  int
	Offset int
}

type ArticleRepo interface {
	// ---- Article 相关 ----

	// Create 创建文章
	Create(ctx context.Context, article *entity.Article) error
	// FindBySlug 根据 slug 查询文章
	FindBySlug(ctx context.Context, slug string) (*entity.Article, error)
	// Update 更新文章
	Update(ctx context.Context, article *entity.Article) error
	// Delete 删除文章
	Delete(ctx context.Context, articleID int64) error

	// List 公开文章列表（支持多条件）
	List(ctx context.Context, query ListArticlesFilter) ([]*entity.Article, int64, error)
	// Feed 关注者文章流
	Feed(ctx context.Context, userID int64, limit, offset int) ([]*entity.Article, int64, error)

	// 这里塞入 tag 和 favorite : Tag / Favorite 是article内部关系
	// ---- Tag 相关 ----

	// GetOrCreateTags 根据 tag 名称获取或创建
	GetOrCreateTags(ctx context.Context, names []string) ([]*entity.Tag, error)
	// ReplaceArticleTags 重置文章标签
	ReplaceArticleTags(ctx context.Context, articleID int64, tags []*entity.Tag) error
	// GetTagsByArticleID 根据 articleID 查 tag
	GetTagsByArticleID(ctx context.Context, articleID int64) ([]*entity.Tag, error)
	ListTags(ctx context.Context) ([]string, error)
	// return map[articleID] []tags
	GetTagsByArticleIDs(ctx context.Context, articleIDs []int64) (map[int64][]string, error)

	// ---- Favorite 相关 ----

	// 是否点赞
	IsFavorited(ctx context.Context, userID, articleID int64) (bool, error)
	// 点赞
	AddFavorite(ctx context.Context, userID, articleID int64) error
	// 取消点赞
	RemoveFavorite(ctx context.Context, userID, articleID int64) error
	CountFavorites(ctx context.Context, articleID int64) (int, error)
}
