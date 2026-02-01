package dto

import "time"

type CreateCommentRequest struct {
	Comment struct {
		Body string `json:"body"`
	} `json:"comment"`
}

type CommentDTO struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    AuthorDTO `json:"author"`
}

type SingleCommentResponse struct {
	Comment CommentDTO `json:"comment"`
}

type MultipleCommentsResponse struct {
	Comments []CommentDTO `json:"comments"`
}
