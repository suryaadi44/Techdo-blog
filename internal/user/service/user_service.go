package service

import (
	"context"
	"database/sql"

	"github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/internal/user/service/impl"
)

type UserServiceApi interface {
	GetUserMiniDetail(ctx context.Context, id int64) (dto.MiniUserDetail, error)
}

func NewUserService(DB *sql.DB) UserServiceApi {
	userRepository := impl.NewUserRepository(DB)

	return impl.UserServiceImpl{
		Repository: userRepository,
	}
}
