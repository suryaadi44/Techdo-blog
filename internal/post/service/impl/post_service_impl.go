package impl

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/codedius/imagekit-go"
	"github.com/suryaadi44/Techdo-blog/internal/post/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type PostServiceImpl struct {
	Repository PostRepositoryImpl
}

func (p PostServiceImpl) AddPost(ctx context.Context, post dto.BlogPostRequest, authorID int64) (int64, error) {
	reservedID, err := p.Repository.ReserveID(ctx)
	if err != nil {
		log.Println("[ERROR] AddPost: Error geting reserved ID-> error:", err)
		return -1, err
	}
	pictureFolder := fmt.Sprintf("/%d", reservedID)

	r := regexp.MustCompile(`\.(\w*)$`)
	extension := r.FindString(post.BannerName)
	bannerName := fmt.Sprintf("%d%s", reservedID, extension)
	bannerUrl, err := utils.UploadImage(ctx, bannerName, post.Banner, pictureFolder)

	r = regexp.MustCompile(`src=\"([^\"]+)\"`)
	matches := r.FindAllStringSubmatch(post.Body, -1)
	for _, v := range matches {
		r := regexp.MustCompile(`image/(\w*)`)
		extension := r.FindAllStringSubmatch(v[1], -1)
		if len(extension) == 0 {
			continue
		}

		pictureName := fmt.Sprintf("%d.%s", reservedID, extension[0][1])
		imgkitResponse, err := utils.UploadImage(ctx, pictureName, v[1], pictureFolder)

		if err == nil {
			post.Body = strings.ReplaceAll(post.Body, v[1], imgkitResponse.URL)
		}
	}

	err = p.Repository.UpdatePost(ctx, post.ToDAO(reservedID, authorID, bannerUrl.URL))
	if err != nil {
		log.Println("[ERROR] AddPost: Error adding post data -> error:", err)
		return -1, err
	}

	err = p.Repository.AddPostCategoryAssoc(ctx, reservedID, post.Category)
	if err != nil {
		log.Println("[ERROR] AddPost: Error adding post category data -> error:", err)
		return -1, err
	}

	return reservedID, nil
}

func (p PostServiceImpl) DeletePost(ctx context.Context, id int64) error {
	err := p.Repository.DeletePost(ctx, id)
	if err != nil {
		log.Println("[ERROR] DeletePost: Error deleting post -> error:", err)
		return err
	}

	imgFodlerPath := fmt.Sprintf("%d/", id)
	err = utils.DeleteFolder(ctx, imgFodlerPath)
	if err != nil {
		log.Println("[ERROR] DeletePost: Error deleting imgkit folder -> error:", err)
		return err
	}

	return nil
}

func (p PostServiceImpl) IncreaseView(ctx context.Context, id int64) error {
	return p.Repository.IncreaseView(ctx, id)
}

func (p PostServiceImpl) GetFullPost(ctx context.Context, id int64) (dto.BlogPostResponse, error) {
	var postDto dto.BlogPostResponse

	post, author, err := p.Repository.GetFullPost(ctx, id)
	if err != nil {
		log.Println("[ERROR] Fetching Full Post with id", id, "-> error:", err)
		return postDto, err
	}

	author.Picture, err = utils.GetPictureUrl(ctx, author.Picture)
	if err != nil {
		log.Println("[ERROR] Error geting user picture url-> error:", err)
		return postDto, err
	}

	categories, err := p.Repository.GetCategoriesFromID(ctx, id)
	if err != nil {
		log.Println("[ERROR] Fetching categories for post with id", id, "-> error:", err)
		return postDto, err
	}

	postDto = dto.NewBlogPostResponse(post, categories, author)
	return postDto, nil
}

func (p PostServiceImpl) GetBriefsBlogPost(ctx context.Context, page int64, limit int64) (dto.BriefsBlogPostResponse, error) {
	var postList dto.BriefsBlogPostResponse
	offset := (page - 1) * limit

	postListEntity, err := p.Repository.GetBriefsBlogPostData(ctx, offset, limit)
	if err != nil {
		log.Println("[ERROR] Fetching list of post -> error:", err)
		return postList, err
	}

	postList = dto.NewBriefsBlogPostResponse(postListEntity)
	return postList, nil
}

func (p PostServiceImpl) GetBriefsBlogPostOfCategories(ctx context.Context, categories string, page int64, limit int64) (dto.BriefsBlogPostResponse, error) {
	var postList dto.BriefsBlogPostResponse
	offset := (page - 1) * limit

	postListEntity, err := p.Repository.GetBriefsBlogPostDataOfCategories(ctx, categories, offset, limit)
	if err != nil {
		log.Println("[ERROR] Fetching list of post -> error:", err)
		return postList, err
	}

	postList = dto.NewBriefsBlogPostResponse(postListEntity)
	return postList, nil
}

func (p PostServiceImpl) GetTopCategoryPost(ctx context.Context) (dto.TopCategoriesWithPost, error) {
	var postData dto.TopCategoriesWithPost

	postRaw, categoryRaw, err := p.Repository.GetTopCategoryPost(ctx)
	if err != nil {
		log.Println("[ERROR] Fetching list of post -> error:", err)
		return postData, err
	}

	postData = dto.NewTopCategoriesAndPost(postRaw, categoryRaw)
	return postData, nil
}

func (p PostServiceImpl) GetEditorsPick(ctx context.Context) (dto.BriefBlogPostResponse, error) {
	var postData dto.BriefBlogPostResponse

	postRaw, err := p.Repository.GetEditorsPick(ctx)
	if err != nil {
		log.Println("[ERROR] Fetching list of post -> error:", err)
		return postData, err
	}

	postData = dto.NewBriefBlogPostResponse(postRaw)
	return postData, nil
}

func (p PostServiceImpl) GetCountListOfPost(ctx context.Context) (int64, error) {
	return p.Repository.CountListOfPost(ctx)
}

func (p PostServiceImpl) GetCountListOfPostInCategories(ctx context.Context, categories string) (int64, error) {
	return p.Repository.CountListOfPostInCategories(ctx, categories)
}

func (p PostServiceImpl) SearchBlogPost(ctx context.Context, q string, page int64, limit int64, dateStart *time.Time, dateEnd *time.Time, category string) (dto.BriefsBlogPostResponse, error) {
	var postList dto.BriefsBlogPostResponse
	offset := (page - 1) * limit

	postListEntity, err := p.Repository.GetBriefsBlogPostFromSearch(ctx, q, offset, limit, dateStart, dateEnd, category)
	if err != nil {
		log.Println("[ERROR] Fetching list of post -> error:", err)
		return postList, err
	}

	postList = dto.NewBriefsBlogPostResponse(postListEntity)
	return postList, nil
}

func (p PostServiceImpl) GetCountOfSearchResult(ctx context.Context, q string, dateStart *time.Time, dateEnd *time.Time, category string) (int64, error) {
	return p.Repository.CountSearchResult(ctx, q, dateStart, dateEnd, category)
}

func (p PostServiceImpl) GetPostAuthorIdFromId(ctx context.Context, postId int64) (int64, error) {
	return p.Repository.GetPostAuthorId(ctx, postId)
}

func (p PostServiceImpl) GetCategoriesFromID(ctx context.Context, id int64) (dto.CategoryList, error) {
	categories, err := p.Repository.GetCategoriesFromID(ctx, id)
	if err != nil {
		log.Println("[ERROR] Fetching categories for post with id", id, "-> error:", err)
		return dto.CategoryList{}, err
	}

	return dto.NewCategoryList(categories), nil
}

func (p PostServiceImpl) GetCategoryList(ctx context.Context) (dto.CategoryList, error) {
	categories, err := p.Repository.GetCategoryList(ctx)
	if err != nil {
		log.Println("[ERROR] Fetching category list -> error:", err)
		return dto.CategoryList{}, err
	}

	return dto.NewCategoryList(categories), nil
}

func (p PostServiceImpl) UploadImage(ctx context.Context, filename string, image interface{}, folderID string) (*imagekit.UploadResponse, error) {
	return utils.UploadImage(ctx, filename, image, folderID)
}

func (p PostServiceImpl) AddComment(ctx context.Context, comment dto.CommentRequest) error {
	return p.Repository.AddComment(ctx, comment.ToDAO())
}

func (p PostServiceImpl) GetComments(ctx context.Context, postID int64) (dto.CommentsResponse, error) {
	var commentResponse dto.CommentsResponse

	comment, user, err := p.Repository.GetPostComments(ctx, postID)
	if err != nil {
		log.Println("[ERROR] Fetching list of comment -> error:", err)
		return commentResponse, err
	}

	commentResponse = dto.NewCommentsResponse(comment, user)

	return commentResponse, nil
}
