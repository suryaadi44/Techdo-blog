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
	IncreaseView(ctx context.Context, id int64) error
	SearchBlogPost(ctx context.Context, q string, page int64, limit int64, dateStart *time.Time, dateEnd *time.Time, category string) (dto.BriefsBlogPostResponse, error)

	AddPost(ctx context.Context, post dto.BlogPostRequest) (int64, error)
	AddComment(ctx context.Context, comment dto.CommentRequest) error

	DeletePost(ctx context.Context, id int64) error
	DeleteComment(ctx context.Context, id int64) error

	EditPost(ctx context.Context, post dto.BlogPostRequest, PostID int64) (int64, error)

	GetRawPost(ctx context.Context, id int64) (dto.RawBlogPostResponse, error)
	GetFullPost(ctx context.Context, id int64) (dto.BlogPostResponse, error)
	GetTopCategoryPost(ctx context.Context) (dto.TopCategoriesWithPost, error)
	GetBriefsBlogPost(ctx context.Context, page int64, limit int64) (dto.BriefsBlogPostResponse, error)
	GetMiniBlogPostsByUser(ctx context.Context, id int64, limit int64) (dto.MiniBlogPostsResponse, error)
	GetBriefsBlogPostOfCategories(ctx context.Context, categories string, page int64, limit int64) (dto.BriefsBlogPostResponse, error)
	GetEditorsPick(ctx context.Context) (dto.BriefBlogPostResponse, error)
	GetPostAuthorIdFromId(ctx context.Context, postId int64) (int64, error)

	GetCountListOfPost(ctx context.Context) (int64, error)
	GetCountOfSearchResult(ctx context.Context, q string, dateStart *time.Time, dateEnd *time.Time, category string) (int64, error)
	GetCountListOfPostInCategories(ctx context.Context, categories string) (int64, error)
	GetUserPostStatisticOfEachCategory(ctx context.Context, id int64) (dto.EachCategoryStats, error)
	GetUserTotalPostCount(ctx context.Context, id int64) (int64, error)
	GetUserTotalCommentCount(ctx context.Context, id int64) (int64, error)

	GetCategoriesFromID(ctx context.Context, id int64) (dto.CategoryList, error)
	GetCategoryList(ctx context.Context) (dto.CategoryList, error)
	GetComments(ctx context.Context, postID int64) (dto.PostCommentsResponse, error)
	GetCommentsByUser(ctx context.Context, uid int64, limit int64) (dto.UserCommentsResponse, error)
	GetCommentAuthorIdFromId(ctx context.Context, id int64) (int64, error)

	UploadImage(ctx context.Context, filename string, image interface{}, folderID string) (*imagekit.UploadResponse, error)

	PickHeaderPost(ctx context.Context, id int64) error
}

func NewPostService(DB *sql.DB) PostServiceApi {
	postRepository := impl.NewPostRepository(DB)

	return impl.PostServiceImpl{
		Repository: postRepository,
	}
}
