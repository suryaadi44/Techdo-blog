package impl

import (
	"context"
	"fmt"
	"log"
	"regexp"

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

	userDetail = dto.NewMiniUserDetailDTO(user)
	return userDetail, nil
}

func (u UserServiceImpl) GetUserDetail(ctx context.Context, id int64) (dto.UserDetailResponse, error) {
	var userDetail dto.UserDetailResponse

	user, err := u.Repository.GetUserDetail(ctx, id)
	if err != nil {
		log.Println("[ERROR] GetUserMiniDetail: Error geting user detail-> error:", err)
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

	response, err := utils.UploadImage(ctx, fileName, picture, "/user")
	if err != nil {
		log.Println("[ERROR] UpdateUserPicture: Error on uploading file-> error:", err)
		return err
	}

	err = u.Repository.UpdateUserPicture(ctx, response.URL, id)
	if err != nil {
		log.Println("[ERROR] UpdateUserPicture: Error updating user detail-> error:", err)
		return err
	}

	return nil
}

func (u UserServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	return u.Repository.DeleteUser(ctx, id)
}
