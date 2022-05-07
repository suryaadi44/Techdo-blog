package service

import (
	"context"
	"database/sql"

	"github.com/codedius/imagekit-go"
	"github.com/suryaadi44/Techdo-blog/internal/post/dto"
	"github.com/suryaadi44/Techdo-blog/internal/post/service/impl"
)

type PostServiceApi interface {
	AddPost(ctx context.Context, post dto.BlogPostRequest, authorID int64) error
	GetFullPost(ctx context.Context, id int64) (dto.BlogPostResponse, error)
	GetBriefsBlogPost(ctx context.Context, page int64, limit int64) (dto.BriefsBlogPostResponse, error)
	GetCategoriesFromID(ctx context.Context, id int64) (dto.CategoryList, error)
	GetCategoryList(ctx context.Context) (dto.CategoryList, error)
	UploadImage(ctx context.Context, filename string, image interface{}, folderID string) (*imagekit.UploadResponse, error)
}

func NewPostService(DB *sql.DB) PostServiceApi {
	postRepository := impl.NewPostRepository(DB)

	return impl.PostServiceImpl{
		Repository: postRepository,
	}
}
