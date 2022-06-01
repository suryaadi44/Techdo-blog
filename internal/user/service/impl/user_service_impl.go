package impl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/suryaadi44/Techdo-blog/internal/global"
	"github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type UserServiceImpl struct {
	Repository     UserRepositoryImpl
	SessionService SessionServiceImpl
}

func (u UserServiceImpl) RegisterUser(ctx context.Context, user dto.SignUpRequest) error {
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	userEntity := entity.User{
		Username: user.Username,
		Password: hash,
	}

	userCount, err := u.Repository.GetUserCount(ctx)
	if err != nil {
		return err
	}

	if userCount == 0 {
		userEntity.Type = 0
	}

	userDetailEntity := entity.UserDetail{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Picture:   global.PICTURE_DEFAULT,
	}

	return u.Repository.NewUser(ctx, userEntity, userDetailEntity)
}

func (u UserServiceImpl) AuthenticateUser(ctx context.Context, user dto.AuthRequest) (entity.Session, error) {
	saved, err := u.Repository.GetUser(ctx, user.Username)
	if err != nil {
		return entity.Session{}, err
	}

	if !utils.CheckPasswordHash(user.Password, saved.Password) {
		return entity.Session{}, errors.New("Inccorect Username or Password")
	}

	log.Println("[Auth] Login :", user.Username, "approved")
	session := entity.Session{
		Token:    uuid.NewString(),
		UID:      saved.UserID,
		ExpireAt: time.Now().Add(time.Duration(global.SESSION_EXPIRE_IN_SECOND) * time.Second),
	}

	err = u.SessionService.NewSession(ctx, session)
	if err != nil {
		return entity.Session{}, err
	}

	return session, nil
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

func (u UserServiceImpl) GetUserCount(ctx context.Context) (int64, error) {
	return u.Repository.GetUserCount(ctx)
}

func (u UserServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	return u.Repository.DeleteUser(ctx, id)
}
