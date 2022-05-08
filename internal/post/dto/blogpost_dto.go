package dto

import (
	"time"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type BlogPostResponse struct {
	PostID     int64
	Author     string
	Categories CategoryList
	Banner     string
	Title      string
	Body       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type BriefBlogPostResponse struct {
	PostID    int64
	Author    string
	Banner    string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BlogPostRequest struct {
	Category   int64
	Banner     []byte
	BannerName string
	Title      string
	Body       string
}

type Category struct {
	CategoryID   int64
	CategoryName string
}

type CategoryList []Category
type BriefsBlogPostResponse []BriefBlogPostResponse

func NewCategory(c entity.Category) Category {
	return Category{
		CategoryID:   c.CategoryID,
		CategoryName: c.CategoryName,
	}
}

func NewCategoryList(c entity.Categories) CategoryList {
	var categoryList CategoryList

	for _, each := range c {
		eachCategory := NewCategory(*each)
		categoryList = append(categoryList, eachCategory)
	}

	return categoryList
}

func NewBlogPostResponse(post entity.BlogPostFull, categories entity.Categories) BlogPostResponse {
	return BlogPostResponse{
		PostID:     post.PostID,
		Author:     post.Author,
		Categories: NewCategoryList(categories),
		Banner:     post.Banner,
		Title:      post.Title,
		Body:       post.Body,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}

func NewBriefBlogPostResponse(post entity.BriefBlogPost) BriefBlogPostResponse {
	return BriefBlogPostResponse{
		PostID:    post.PostID,
		Author:    post.Author,
		Banner:    post.Banner,
		Title:     post.Title,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

}

func NewBriefsBlogPostResponse(posts entity.BriefsBlogPost) BriefsBlogPostResponse {
	var postList BriefsBlogPostResponse

	for _, each := range posts {
		eachPost := NewBriefBlogPostResponse(*each)
		postList = append(postList, eachPost)
	}

	return postList
}

func (b *BlogPostRequest) ToDAO(PostID int64, AuthorID int64, BannerURL string) entity.BlogPost {
	return entity.BlogPost{
		PostID:   PostID,
		AuthorID: AuthorID,
		Banner:   BannerURL,
		Title:    b.Title,
		Body:     b.Body,
	}
}
