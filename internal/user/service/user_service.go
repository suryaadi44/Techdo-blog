package service

import (
	"context"
	"database/sql"

	"github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/internal/user/service/impl"
)

type UserServiceApi interface {
	UpdateUserDetail(ctx context.Context, user dto.UserDetailRequest) error
	UpdateUserPicture(ctx context.Context, picture []byte, Filename string, id int64) error

	GetUserMiniDetail(ctx context.Context, id int64) (dto.MiniUserDetailResponse, error)
	GetUserDetail(ctx context.Context, id int64) (dto.UserDetailResponse, error)

	DeleteUser(ctx context.Context, id int64) error
}

func NewUserService(DB *sql.DB) UserServiceApi {
	userRepository := impl.NewUserRepository(DB)

	return impl.UserServiceImpl{
		Repository: userRepository,
	}
}
