package service

import (
	"context"
	"database/sql"

	"github.com/codedius/imagekit-go"
	"github.com/suryaadi44/Techdo-blog/internal/post/dto"
	"github.com/suryaadi44/Techdo-blog/internal/post/service/impl"
)

type PostServiceApi interface {
	GetFullPost(ctx context.Context, id int) (dto.BlogPost, error)
	GetCategoriesFromID(ctx context.Context, id int) (dto.CategoryList, error)
	GetCategoryList(ctx context.Context) (dto.CategoryList, error)
	UploadImage(ctx context.Context, filename string, image []byte) (*imagekit.UploadResponse, error)
}

func NewPostService(DB *sql.DB) PostServiceApi {
	postRepository := impl.NewPostRepository(DB)

	return impl.PostServiceImpl{
		Repository: postRepository,
	}
}
