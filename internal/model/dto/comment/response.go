package comment

import "time"

type CommentDTO struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    struct {
		Username  string `json:"username"`
		Bio       string `json:"bio"`
		Image     string `json:"image"`
		Following bool   `json:"following"`
	} `json:"author"`
}
type SingleCommentResponse struct {
	Comment CommentDTO `json:"comment"`
}
type MultipleCommentResponse struct {
	Comments []CommentDTO `json:"comments"`
}

type ListofTagsResponse struct {
	Tags []string `json:"tags"`
}
