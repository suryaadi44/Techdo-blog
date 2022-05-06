package dto

import (
	"time"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type BlogPost struct {
	PostID     int64
	Author     string
	Categories CategoryList
	Banner     string
	Title      string
	Body       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Category struct {
	CategoryID   int64
	CategoryName string
}

type CategoryList []Category

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

func NewBlogPost(post entity.BlogPostFull, categories entity.Categories) BlogPost {
	return BlogPost{
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
