package impl

import (
	"context"
	"log"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type SessionServiceImpl struct {
	Repository SessionRepositoryImpl
}

func (s SessionServiceImpl) NewSession(ctx context.Context, user entity.Session) error {
	return s.Repository.NewSession(ctx, user)
}

func (s SessionServiceImpl) GetSession(ctx context.Context, token string) (entity.SessionDetail, error) {
	session, err := s.Repository.GetSession(ctx, token)
	if err != nil {
		log.Println("[Session] GetSession error on getting session", err.Error())
		return session, err
	}

	return session, nil
}

func (s SessionServiceImpl) DeleteSession(ctx context.Context, token string) error {
	return s.Repository.DeleteSession(ctx, token)
}
