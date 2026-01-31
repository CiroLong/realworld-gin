package dto

import (
	"github/CiroLong/realworld-gin/internal/model/entity"
	"time"
)

type AuthorDTO struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type ArticleDTO struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"` //TODO: check所有返回article，需要填写这个
	FavoritesCount int       `json:"favoritesCount"`
	Author         AuthorDTO `json:"author"`
}

type ArticleResponse struct {
	Article ArticleDTO `json:"article"`
}

func NewArticleResponse(article *entity.Article, tags []string, author AuthorDTO, favorited bool) *ArticleResponse {
	return &ArticleResponse{
		Article: ArticleDTO{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        tags,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      article.UpdatedAt,
			Favorited:      favorited,
			FavoritesCount: article.FavoritesCount,
			Author:         author,
		},
	}
}

// 注意这里不返回Body
type ArticleWithoutBodyDTO struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         AuthorDTO `json:"author"`
}
type MultipleArticlesResponse struct {
	Articles      []ArticleWithoutBodyDTO `json:"articles"`
	ArticlesCount int                     `json:"articlesCount"`
}
