package dto

import (
	"html/template"
	"regexp"

	UserDto "github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type RawBlogPostResponse struct {
	PostID     int64         `json:"postId"`
	AuthorID   int64         `json:"authorId"`
	Banner     string        `json:"banner"`
	Title      string        `json:"title"`
	Categories CategoryList  `json:"categories"`
	Body       template.HTML `json:"body"`
	CreatedAt  string        `json:"createdAt"`
	UpdatedAt  string        `json:"updatedAt"`
}

type BlogPostResponse struct {
	PostID       int64                      `json:"postId"`
	Author       UserDto.UserDetailResponse `json:"author"`
	Categories   CategoryList               `json:"categories"`
	Banner       string                     `json:"banner"`
	Title        string                     `json:"title"`
	Body         template.HTML              `json:"body"`
	ViewCount    int64                      `json:"viewCount"`
	CommentCount int64                      `json:"commentCount"`
	CreatedAt    string                     `json:"createdAt"`
	UpdatedAt    string                     `json:"updatedAt"`
}

type BriefBlogPostResponse struct {
	PostID       int64  `json:"postId"`
	Author       string `json:"author"`
	Banner       string `json:"banner"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	ViewCount    int64  `json:"viewCount"`
	CommentCount int64  `json:"commentCount"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type MiniBlogPostResponse struct {
	Index     int64          `json:"index"`
	PostID    dto.NullInt64  `json:"postId"`
	Title     dto.NullString `json:"title"`
	CreatedAt dto.NullTime   `json:"createdAt"`
	Category  dto.NullString `json:"category"`
}

type BlogPostRequest struct {
	AuthorID   int64  `json:"authorId"`
	Category   int64  `json:"category"`
	Banner     []byte `json:"banner"`
	BannerName string `json:"bannerName"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}

type Category struct {
	CategoryID   int64  `json:"categoryID"`
	CategoryName string `json:"categoryName"`
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

func NewRawBlogPostResponse(post entity.BlogPost, categories entity.Categories) RawBlogPostResponse {
	return RawBlogPostResponse{
		PostID:     post.PostID,
		AuthorID:   post.AuthorID,
		Banner:     post.Banner,
		Title:      post.Title,
		Categories: NewCategoryList(categories),
		Body:       template.HTML(post.Body),
		CreatedAt:  post.CreatedAt.Format("Jan 02, 2006"),
		UpdatedAt:  post.UpdatedAt.Format("Jan 02, 2006"),
	}
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

func NewMiniBlogPostResponse(post entity.PostTitleWithCategory, index int64) MiniBlogPostResponse {
	return MiniBlogPostResponse{
		Index:     index,
		PostID:    post.PostID,
		Title:     post.Title,
		CreatedAt: post.CreatedAt,
		Category:  post.Category,
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

func NewMiniBlogPostsResponse(posts entity.PostsTitleWithCategory) MiniBlogPostsResponse {
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

func (b *BlogPostRequest) ToDAO(PostID int64, BannerURL string) entity.BlogPost {
	return entity.BlogPost{
		PostID:   PostID,
		AuthorID: b.AuthorID,
		Banner:   BannerURL,
		Title:    b.Title,
		Body:     b.Body,
	}
}
