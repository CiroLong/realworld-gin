package service

import (
	"context"
	"errors"
	"github/CiroLong/realworld-gin/internal/model/dto"
)

var ErrPermissionDenied = errors.New("permission denied")

type ArticleService interface {
	CreateArticle(ctx context.Context, authorID int64, req *dto.CreateArticleRequest) (*dto.ArticleResponse, error)
	GetArticle(ctx context.Context, slug string) (*dto.ArticleResponse, error)
	UpdateArticle(ctx context.Context, slug string, userID int64, req *dto.UpdateArticleRequest) (*dto.ArticleResponse, error)
	DeleteArticle(ctx context.Context, slug string, userID int64) error
	FavoriteArticle(ctx context.Context, slug string, userID int64) (*dto.ArticleResponse, error)
	UnfavoriteArticle(ctx context.Context, slug string, userID int64) (*dto.ArticleResponse, error)

	ListArticles(ctx context.Context, tag string, author string, favorited string, userID int64, limit int, offset int) (*dto.MultipleArticlesResponse, error)
	FeedArticles(ctx context.Context, userID int64, limit int, offset int) (*dto.MultipleArticlesResponse, error)

	ListTags(ctx context.Context) ([]string, error)
}
