package impl

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/suryaadi44/Techdo-blog/internal/global"
	"github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type UserServiceImpl struct {
	Repository UserRepositoryImpl
}

func (u UserServiceImpl) GetUserMiniDetail(ctx context.Context, id int64) (dto.MiniUserDetailResponse, error) {
	var userDetail dto.MiniUserDetailResponse

	user, err := u.Repository.GetUserMiniDetail(ctx, id)
	if err != nil {
		log.Println("[ERROR] GetUserMiniDetail: Error geting user detail-> error:", err)
		return userDetail, err
	}

	user.Picture, err = utils.GetPictureUrl(ctx, user.Picture)
	if err != nil {
		log.Println("[ERROR] GetUserMiniDetail: Error geting user picture url-> error:", err)
		return userDetail, err
	}

	userDetail = dto.NewMiniUserDetailDTO(user)
	return userDetail, nil
}

func (u UserServiceImpl) GetUserDetail(ctx context.Context, id int64) (dto.UserDetailResponse, error) {
	var userDetail dto.UserDetailResponse

	user, err := u.Repository.GetUserDetail(ctx, id)
	if err != nil {
		log.Println("[ERROR] GetUserDetail: Error geting user detail-> error:", err)
		return userDetail, err
	}

	user.Picture, err = utils.GetPictureUrl(ctx, user.Picture)
	if err != nil {
		log.Println("[ERROR] GetUserDetail: Error geting user picture url-> error:", err)
		return userDetail, err
	}

	userDetail = dto.NewUserDetailDTO(user)
	return userDetail, nil
}

func (u UserServiceImpl) UpdateUserDetail(ctx context.Context, user dto.UserDetailRequest) error {
	userEntity := dto.NewUserDetailDAO(user)

	err := u.Repository.UpdateUserDetail(ctx, userEntity)
	if err != nil {
		log.Println("[ERROR] UpdateUserDetail: Error updating user detail-> error:", err)
		return err
	}

	return nil
}

func (u UserServiceImpl) UpdateUserPicture(ctx context.Context, picture []byte, tempName string, id int64) error {
	r := regexp.MustCompile(`\.(\w*)$`)
	extension := r.FindString(tempName)
	fileName := fmt.Sprintf("%d%s", id, extension)

	oldID, err := u.Repository.GetUserPictureID(ctx, id)
	if err != nil {
		log.Println("[ERROR] UpdateUserPicture: Error on getting old picture id-> error:", err)
		return err
	}

	if oldID != global.PICTURE_DEFAULT {
		if err := utils.DeleteImage(ctx, oldID); err != nil {
			log.Println("[ERROR] UpdateUserPicture: Error on deleting old picture-> error:", err)
			return err
		}
	}

	response, err := utils.UploadImage(ctx, fileName, picture, "/user")
	if err != nil {
		log.Println("[ERROR] UpdateUserPicture: Error on uploading file-> error:", err)
		return err
	}

	err = u.Repository.UpdateUserPicture(ctx, response.FileID, id)
	if err != nil {
		log.Println("[ERROR] UpdateUserPicture: Error updating user detail-> error:", err)
		return err
	}

	return nil
}

func (u UserServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	return u.Repository.DeleteUser(ctx, id)
}

func (u UserServiceImpl) GetUserTotalPostCount(ctx context.Context, id int64) (int64, error) {
	total, err := u.Repository.GetUserTotalPostCount(ctx, id)
	if err != nil {
		log.Println("[ERROR] GetUserTotalPostCount: Error getting count of user total post-> error:", err)
		return 0, err
	}

	return total, nil
}

func (u UserServiceImpl) GetUserTotalCommentCount(ctx context.Context, id int64) (int64, error) {
	total, err := u.Repository.GetUserTotalCommentCount(ctx, id)
	if err != nil {
		log.Println("[ERROR] GetUserTotalPostCount: Error getting count of user total post-> error:", err)
		return 0, err
	}

	return total, nil
}
