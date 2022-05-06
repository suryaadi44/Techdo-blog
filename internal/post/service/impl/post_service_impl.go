package impl

import (
	"context"
	"log"

	"github.com/codedius/imagekit-go"
	"github.com/suryaadi44/Techdo-blog/internal/post/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type PostServiceImpl struct {
	Repository PostRepositoryImpl
}

func (p PostServiceImpl) GetFullPost(ctx context.Context, id int) (dto.BlogPost, error) {
	post, err := p.Repository.GetFullPost(ctx, id)
	if err != nil {
		log.Println("ERROR Fetching Full Post with id", id, "-> error:", err)
		return dto.BlogPost{}, err
	}

	categories, err := p.Repository.GetCategoriesFromID(ctx, id)
	if err != nil {
		log.Println("ERROR Fetching categories for post with id", id, "-> error:", err)
		return dto.BlogPost{}, err
	}

	return dto.NewBlogPost(post, categories), nil
}

func (p PostServiceImpl) GetCategoriesFromID(ctx context.Context, id int) (dto.CategoryList, error) {
	categories, err := p.Repository.GetCategoriesFromID(ctx, id)
	if err != nil {
		log.Println("ERROR Fetching categories for post with id", id, "-> error:", err)
		return dto.CategoryList{}, err
	}

	return dto.NewCategoryList(categories), nil
}

func (p PostServiceImpl) GetCategoryList(ctx context.Context) (dto.CategoryList, error) {
	categories, err := p.Repository.GetCategoryList(ctx)
	if err != nil {
		log.Println("ERROR Fetching category list -> error:", err)
		return dto.CategoryList{}, err
	}

	return dto.NewCategoryList(categories), nil
}

func (p PostServiceImpl) UploadImage(ctx context.Context, image dto.Image) (*imagekit.UploadResponse, error) {
	folder := "/posts"
	return utils.UploadImage(ctx, image.FileName, image.Data, folder)
}
