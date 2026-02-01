package service

import (
	"context"
	"github/CiroLong/realworld-gin/internal/model/dto"
	"github/CiroLong/realworld-gin/internal/model/entity"
	"github/CiroLong/realworld-gin/internal/pkg/common"
	"github/CiroLong/realworld-gin/internal/repository"
)

type commentService struct {
	commentRepo repository.CommentRepo
	articleRepo repository.ArticleRepo
	userRepo    repository.UserRepo
}

func NewCommentService(commentRepo repository.CommentRepo,
	articleRepo repository.ArticleRepo,
	userRepo repository.UserRepo,
) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

func (c commentService) CreateComment(ctx context.Context, userID int64, slug string, req *dto.CreateCommentRequest) (*dto.SingleCommentResponse, error) {
	// 1. 查文章
	article, err := c.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// 2. 创建 Comment entity
	comment := &entity.Comment{
		Body:      req.Comment.Body,
		ArticleID: article.ID,
		AuthorID:  userID,
	}

	if err = c.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// 3. 查作者
	author, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 4. 组装DTO ,following（自己一定是 false）
	authorDTO := dto.AuthorDTO{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: false,
	}

	commentDTO := dto.CommentDTO{
		ID:        comment.ID,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Author:    authorDTO,
	}

	return &dto.SingleCommentResponse{
		Comment: commentDTO,
	}, nil
}

// userID == 0 时不用查following
func (c commentService) GetComments(ctx context.Context, slug string, userID int64) (*dto.MultipleCommentsResponse, error) {
	// 1. 查文章
	article, err := c.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// 2. 查评论
	comments, err := c.commentRepo.ListByArticle(ctx, article.ID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.CommentDTO, 0, len(comments))

	for _, comment := range comments {
		// 3. 查作者
		// TODO: 批量查 user
		author, err := c.userRepo.FindByID(ctx, comment.AuthorID)
		if err != nil {
			return nil, err
		}

		// TODO: 优化：批量查following关系
		following := false
		if userID != 0 && userID != author.ID {
			following, err = c.userRepo.IsFollowing(ctx, userID, author.ID)
			if err != nil {
				return nil, err
			}
		}

		authorDTO := dto.AuthorDTO{
			Username:  author.Username,
			Bio:       author.Bio,
			Image:     author.Image,
			Following: following,
		}

		result = append(result, dto.CommentDTO{
			ID:        comment.ID,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Author:    authorDTO,
		})
	}

	return &dto.MultipleCommentsResponse{
		Comments: result,
	}, nil
}

func (c commentService) DeleteComment(ctx context.Context, userID int64, commentID int64) error {
	// 1. 查评论
	comment, err := c.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}

	// 2. 权限校验
	if comment.AuthorID != userID {
		return common.ErrPermissionDenied
	}

	// 3. 删除
	return c.commentRepo.Delete(ctx, commentID)
}
