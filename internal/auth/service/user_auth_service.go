package service

import (
	"context"
	"database/sql"

	"github.com/suryaadi44/Techdo-blog/internal/auth/dto"
	"github.com/suryaadi44/Techdo-blog/internal/auth/service/impl"
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type UserAuthServiceApi interface {
	RegisterUser(ctx context.Context, user dto.AuthRequest) error
	AuthenticateUser(ctx context.Context, user dto.AuthRequest) (entity.Session, error)
}

func NewUserAuthService(DB *sql.DB, SessionService impl.SessionServiceImpl) UserAuthServiceApi {
	authRepository := impl.NewUserAuthRepository(DB)

	return impl.UserAuthServiceImpl{
		Repository:     authRepository,
		SessionService: SessionService,
	}
}
