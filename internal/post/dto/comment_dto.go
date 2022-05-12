package dto

import (
	"time"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type CommentRequest struct {
	PostID      int64  `json:"post_id"`
	UserID      int64  `json:"uid"`
	CommentBody string `json:"commentBody"`
}

type CommentResponse struct {
	UserID      int64     `json:"uid"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	UserPicture string    `json:"userPic"`
	CommentBody string    `json:"commentBody"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CommentsResponse []*CommentResponse

func (c *CommentRequest) ToDAO() entity.Comment {
	return entity.Comment{
		PostID:      c.PostID,
		UserID:      c.UserID,
		CommentBody: c.CommentBody,
	}
}

func NewCommentResponse(comment entity.Comment, user entity.MiniUserDetail) CommentResponse {
	return CommentResponse{
		UserID:      comment.UserID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		UserPicture: user.Picture,
		CommentBody: comment.CommentBody,
		CreatedAt:   comment.CreatedAt,
		UpdatedAt:   comment.UpdatedAt,
	}
}

func NewCommentsResponse(comments entity.Comments, users entity.MiniUsersDetail) CommentsResponse {
	var response CommentsResponse

	for idx := range comments {
		eachComment := comments[idx]
		eachUsers := users[idx]

		comment := NewCommentResponse(*eachComment, *eachUsers)
		response = append(response, &comment)
	}

	return response
}
