package dto

import (
	"html/template"
	"regexp"

	UserDto "github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type BlogPostResponse struct {
	PostID       int64
	Author       UserDto.UserDetailResponse
	Categories   CategoryList
	Banner       string
	Title        string
	Body         template.HTML
	ViewCount    int64
	CommentCount int64
	CreatedAt    string
	UpdatedAt    string
}

type BriefBlogPostResponse struct {
	PostID       int64
	Author       string
	Banner       string
	Title        string
	Body         string
	ViewCount    int64
	CommentCount int64
	CreatedAt    string
	UpdatedAt    string
}

type MiniBlogPostResponse struct {
	Index     int64
	PostID    int64
	Title     string
	CreatedAt string
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
type MiniBlogPostsResponse []MiniBlogPostResponse
type TopCategoriesWithPost map[string]BriefsBlogPostResponse

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
		PostID:       post.PostID,
		Author:       UserDto.NewUserDetailDTO(author),
		Categories:   NewCategoryList(categories),
		Banner:       post.Banner,
		Title:        post.Title,
		Body:         template.HTML(post.Body),
		ViewCount:    post.ViewCount,
		CommentCount: post.CommentCount,
		CreatedAt:    post.CreatedAt.Format("Jan 02, 2006"),
		UpdatedAt:    post.UpdatedAt.Format("Jan 02, 2006"),
	}
}

func NewBriefBlogPostResponse(post entity.BriefBlogPost) BriefBlogPostResponse {
	r := regexp.MustCompile(`<[^>]*>`)
	body := utils.Truncate(r.ReplaceAllString(post.Body, ""), 150)

	return BriefBlogPostResponse{
		PostID:       post.PostID,
		Author:       post.Author,
		Banner:       post.Banner,
		Title:        post.Title,
		Body:         body,
		ViewCount:    post.ViewCount,
		CommentCount: post.CommentCount,
		CreatedAt:    post.CreatedAt.Format("Jan 02, 2006"),
		UpdatedAt:    post.UpdatedAt.Format("Jan 02, 2006"),
	}
}

func NewMiniBlogPostResponse(post entity.BriefBlogPost, index int64) MiniBlogPostResponse {
	return MiniBlogPostResponse{
		Index:     index,
		PostID:    post.PostID,
		Title:     post.Title,
		CreatedAt: post.CreatedAt.Format("Jan 02, 2006"),
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

func NewMiniBlogPostsResponse(posts entity.BriefsBlogPost) MiniBlogPostsResponse {
	var postList MiniBlogPostsResponse

	for idx, each := range posts {
		eachPost := NewMiniBlogPostResponse(*each, int64(idx+1))
		postList = append(postList, eachPost)
	}

	return postList
}

func NewTopCategoriesAndPost(posts entity.BriefsBlogPost, categories entity.Categories) TopCategoriesWithPost {
	var topCategoriesAndPost = make(TopCategoriesWithPost)

	for idx, each := range categories {
		postData := NewBriefBlogPostResponse(*posts[idx])
		topCategoriesAndPost[each.CategoryName] = append(topCategoriesAndPost[each.CategoryName], postData)
	}

	return topCategoriesAndPost
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
