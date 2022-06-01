package entity

import (
	"time"

	"github.com/suryaadi44/Techdo-blog/pkg/dto"
)

type BlogPost struct {
	PostID       int64     `db:"post_id"`
	AuthorID     int64     `db:"author_id"`
	Banner       string    `db:"banner"`
	Title        string    `db:"title"`
	Body         string    `db:"body"`
	ViewCount    int64     `db:"view_count"`
	CommentCount int64     `db:"comment_count"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type BriefBlogPost struct {
	PostID       int64     `db:"post_id"`
	Banner       string    `db:"banner"`
	Title        string    `db:"title"`
	Body         string    `db:"body"`
	ViewCount    int64     `db:"view_count"`
	CommentCount int64     `db:"comment_count"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Author       string    `db:"author"`
}
type PostTitleWithCategory struct {
	PostID    dto.NullInt64  `db:"post_id"`
	Title     dto.NullString `db:"title"`
	CreatedAt dto.NullTime   `db:"created_at"`
	Category  dto.NullString `db:"category_name"`
}

type UserPostStatisticByCategory struct {
	Category  dto.NullString `db:"category_name"`
	TotalPost int64          `db:"TotalPost"`
	TotalView int64          `db:"TotalView"`
}

type Category struct {
	CategoryID   int64  `db:"category_id"`
	CategoryName string `db:"category_name"`
}

type Comment struct {
	CommentID   int64     `db:"comment_id"`
	PostID      int64     `db:"post_id"`
	UserID      int64     `db:"uid"`
	CommentBody string    `db:"comment_body"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Categories []*Category
type Comments []*Comment
type BriefsBlogPost []*BriefBlogPost
type PostsTitleWithCategory []*PostTitleWithCategory
type ListOfUserPostStatisticByCategory []*UserPostStatisticByCategory
