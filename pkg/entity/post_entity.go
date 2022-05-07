package entity

import "time"

type BlogPost struct {
	PostID    int64     `db:"post_id"`
	AuthorID  int64     `db:"author_id"`
	Banner    string    `db:"banner"`
	Title     string    `db:"title"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type BlogPostFull struct {
	PostID    int64     `db:"post_id"`
	Banner    string    `db:"banner"`
	Title     string    `db:"title"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Author    string    `db:"author"`
	Picture   string    `db:"picture"`
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
