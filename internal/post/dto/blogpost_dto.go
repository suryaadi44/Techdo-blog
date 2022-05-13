package dto

import (
	"html/template"
	"regexp"

	UserDto "github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type BlogPostResponse struct {
	PostID     int64
	Author     UserDto.UserDetail
	Categories CategoryList
	Banner     string
	Title      string
	Body       template.HTML
	CreatedAt  string
	UpdatedAt  string
}

type BriefBlogPostResponse struct {
	PostID    int64
	Author    string
	Banner    string
	Title     string
	Body      string
	CreatedAt string
	UpdatedAt string
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

func NewBlogPostResponse(post entity.BlogPost, categories entity.Categories, author entity.UserDetail) BlogPostResponse {
	return BlogPostResponse{
		PostID:     post.PostID,
		Author:     UserDto.NewUserDetailDTO(author),
		Categories: NewCategoryList(categories),
		Banner:     post.Banner,
		Title:      post.Title,
		Body:       template.HTML(post.Body),
		CreatedAt:  post.CreatedAt.Format("02 Jan 2006"),
		UpdatedAt:  post.UpdatedAt.Format("02 Jan 2006"),
	}
}

func NewBriefBlogPostResponse(post entity.BriefBlogPost) BriefBlogPostResponse {
	r := regexp.MustCompile(`<[^>]*>`)
	body := utils.Truncate(r.ReplaceAllString(post.Body, ""), 230)

	return BriefBlogPostResponse{
		PostID:    post.PostID,
		Author:    post.Author,
		Banner:    post.Banner,
		Title:     post.Title,
		Body:      body,
		CreatedAt: post.CreatedAt.Format("02 Jan 2006"),
		UpdatedAt: post.UpdatedAt.Format("02 Jan 2006"),
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
