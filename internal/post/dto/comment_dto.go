package dto

import (
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type CommentRequest struct {
	PostID      int64  `json:"post_id"`
	UserID      int64  `json:"uid"`
	CommentBody string `json:"commentBody"`
}

type PostCommentResponse struct {
	UserID      int64  `json:"uid"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	UserPicture string `json:"userPic"`
	CommentBody string `json:"commentBody"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type UserCommentResponse struct {
	Index       int64  `json:"index"`
	CommentID   int64  `json:"commentID"`
	PostID      int64  `json:"postID"`
	PostTitle   string `json:"postTitle"`
	CommentBody string `json:"commentBody"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type PostCommentsResponse []*PostCommentResponse
type UserCommentsResponse []*UserCommentResponse

func (c *CommentRequest) ToDAO() entity.Comment {
	return entity.Comment{
		PostID:      c.PostID,
		UserID:      c.UserID,
		CommentBody: c.CommentBody,
	}
}

func NewPostCommentResponse(comment entity.Comment, user entity.MiniUserDetail) PostCommentResponse {
	return PostCommentResponse{
		UserID:      comment.UserID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		UserPicture: user.Picture,
		CommentBody: comment.CommentBody,
		CreatedAt:   comment.CreatedAt.Format("15:04, Jan 02, 2006"),
		UpdatedAt:   comment.UpdatedAt.Format("15:04, Jan 02, 2006"),
	}
}

func NewPostCommentsResponse(comments entity.Comments, users entity.MiniUsersDetail) PostCommentsResponse {
	var response PostCommentsResponse

	for idx := range comments {
		eachComment := comments[idx]
		eachUsers := users[idx]

		comment := NewPostCommentResponse(*eachComment, *eachUsers)
		response = append(response, &comment)
	}

	return response
}

func NewUserCommentResponse(comment entity.Comment, post entity.BriefBlogPost, index int64) UserCommentResponse {
	return UserCommentResponse{
		Index:       index,
		CommentID:   comment.CommentID,
		PostID:      post.PostID,
		PostTitle:   post.Title,
		CommentBody: comment.CommentBody,
		CreatedAt:   comment.CreatedAt.Format("15:04, Jan 02, 2006"),
		UpdatedAt:   comment.UpdatedAt.Format("15:04, Jan 02, 2006"),
	}
}

func NewUserCommentsResponse(comments entity.Comments, posts entity.BriefsBlogPost) UserCommentsResponse {
	var response UserCommentsResponse

	for idx := range comments {
		comment := NewUserCommentResponse(*comments[idx], *posts[idx], int64(idx+1))
		response = append(response, &comment)
	}

	return response
}
