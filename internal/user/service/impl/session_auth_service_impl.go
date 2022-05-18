package impl

import (
	"context"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type SessionServiceImpl struct {
	Repository SessionRepositoryImpl
}

func (s SessionServiceImpl) NewSession(ctx context.Context, user entity.Session) error {
	return s.Repository.NewSession(ctx, user)
}

func (s SessionServiceImpl) GetSession(ctx context.Context, token string) (entity.SessionDetail, error) {
	return s.Repository.GetSession(ctx, token)
}

func (s SessionServiceImpl) DeleteSession(ctx context.Context, token string) error {
	return s.Repository.DeleteSession(ctx, token)
}
