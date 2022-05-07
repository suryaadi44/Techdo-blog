package service

import (
	"context"
	"database/sql"

	"github.com/suryaadi44/Techdo-blog/internal/auth/service/impl"
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type SessionServiceApi interface {
	NewSession(ctx context.Context, user entity.Session) error
	GetSession(ctx context.Context, token string) (entity.SessionDetail, error)
	DeleteSession(ctx context.Context, token string) error
}

func NewSessionAuthService(DB *sql.DB) impl.SessionServiceImpl {
	sessionRepository := impl.NewSessionRepository(DB)

	return impl.SessionServiceImpl{
		Repository: sessionRepository,
	}
}
