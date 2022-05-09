package impl

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

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

func (p PostServiceImpl) GetFullPost(ctx context.Context, id int64) (dto.BlogPostResponse, error) {
	post, err := p.Repository.GetFullPost(ctx, id)
	if err != nil {
		log.Println("[ERROR] Fetching Full Post with id", id, "-> error:", err)
		return dto.BlogPostResponse{}, err
	}

	categories, err := p.Repository.GetCategoriesFromID(ctx, id)
	if err != nil {
		log.Println("[ERROR] Fetching categories for post with id", id, "-> error:", err)
		return dto.BlogPostResponse{}, err
	}

	return dto.NewBlogPostResponse(post, categories), nil
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
