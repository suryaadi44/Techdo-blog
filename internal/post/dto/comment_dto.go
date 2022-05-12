package dto

import "github.com/suryaadi44/Techdo-blog/pkg/entity"

type CommentRequest struct {
	PostID      int64  `json:"post_id"`
	UserID      int64  `json:"uid"`
	CommentBody string `json:"comment_body"`
}

func (c *CommentRequest) ToDAO() entity.Comment {
	return entity.Comment{
		PostID:      c.PostID,
		UserID:      c.UserID,
		CommentBody: c.CommentBody,
	}
}
