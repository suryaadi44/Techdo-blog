package impl

import (
	"context"
	"log"

	"github.com/suryaadi44/Techdo-blog/internal/user/dto"
)

type UserServiceImpl struct {
	Repository UserRepositoryImpl
}

func (u UserServiceImpl) GetUserMiniDetail(ctx context.Context, id int64) (dto.MiniUserDetail, error) {
	var userDetail dto.MiniUserDetail

	user, err := u.Repository.GetUserDetail(ctx, id)
	if err != nil {
		log.Println("[ERROR] GetUserMiniDetail: Error geting user detail-> error:", err)
		return userDetail, err
	}

	userDetail = dto.NewMiniUserDetailDTO(user)
	return userDetail, nil
}
