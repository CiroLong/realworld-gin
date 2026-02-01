package service

import (
	"context"
	"github/CiroLong/realworld-gin/internal/model/dto"
)

type CommentService interface {
	CreateComment(ctx context.Context, userID int64, slug string, req *dto.CreateCommentRequest) (*dto.SingleCommentResponse, error)

	GetComments(ctx context.Context, slug string, userID int64) (*dto.MultipleCommentsResponse, error)

	DeleteComment(ctx context.Context, userID int64, commentID int64) error
}
