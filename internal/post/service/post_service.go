package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/codedius/imagekit-go"
	"github.com/suryaadi44/Techdo-blog/internal/post/dto"
	"github.com/suryaadi44/Techdo-blog/internal/post/service/impl"
)

type PostServiceApi interface {
	SearchBlogPost(ctx context.Context, q string, page int64, limit int64, dateStart *time.Time, dateEnd *time.Time) (dto.BriefsBlogPostResponse, error)

	AddPost(ctx context.Context, post dto.BlogPostRequest, authorID int64) (int64, error)
	AddComment(ctx context.Context, comment dto.CommentRequest) error

	DeletePost(ctx context.Context, id int64) error

	GetFullPost(ctx context.Context, id int64) (dto.BlogPostResponse, error)
	GetBriefsBlogPost(ctx context.Context, page int64, limit int64) (dto.BriefsBlogPostResponse, error)
	GetPostAuthorIdFromId(ctx context.Context, postId int64) (int64, error)
	GetCategoriesFromID(ctx context.Context, id int64) (dto.CategoryList, error)
	GetCategoryList(ctx context.Context) (dto.CategoryList, error)
	GetComments(ctx context.Context, postID int64, page int64, limit int64) (dto.CommentsResponse, error)

	UploadImage(ctx context.Context, filename string, image interface{}, folderID string) (*imagekit.UploadResponse, error)
}

func NewPostService(DB *sql.DB) PostServiceApi {
	postRepository := impl.NewPostRepository(DB)

	return impl.PostServiceImpl{
		Repository: postRepository,
	}
}
