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

func (p PostServiceImpl) AddPost(ctx context.Context, post dto.BlogPostRequest, authorID int64) error {
	reservedID, err := p.Repository.ReserveID(ctx)
	if err != nil {
		log.Println("[ERROR] AddPost: Error geting reserved ID-> error:", err)
		return err
	}

	r := regexp.MustCompile(`src="([^"]+)"`)
	matches := r.FindAllStringSubmatch(post.Body, -1)
	pictureFolder := fmt.Sprintf("/%d", reservedID)
	for _, v := range matches {
		r := regexp.MustCompile(`image/(\w*)`)
		extension := r.FindAllStringSubmatch(v[1], -1)[0][1]

		pictureName := fmt.Sprintf("%d.%s", reservedID, extension)
		imgkitResponse, err := utils.UploadImage(ctx, pictureName, v[1], pictureFolder)

		if err == nil {
			post.Body = strings.ReplaceAll(post.Body, v[1], imgkitResponse.URL)
		}
	}

	err = p.Repository.UpdatePost(ctx, post.ToDAO(reservedID, authorID))
	if err != nil {
		log.Println("[ERROR] AddPost: Error adding post data -> error:", err)
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

func (p PostServiceImpl) UploadImage(ctx context.Context, filename string, image string, folderID string) (*imagekit.UploadResponse, error) {

	return utils.UploadImage(ctx, filename, image, folderID)
}
