package impl

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/suryaadi44/Techdo-blog/internal/auth/dto"
	"github.com/suryaadi44/Techdo-blog/internal/global"
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type UserAuthServiceImpl struct {
	Repository     UserAuthRepositoryImpl
	SessionService SessionServiceImpl
}

func (u UserAuthServiceImpl) RegisterUser(ctx context.Context, user dto.AuthRequest) error {
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	data := entity.User{
		Username: user.Username,
		Password: hash,
	}

	err = u.Repository.NewUser(ctx, data)

	return err
}

func (u UserAuthServiceImpl) AuthenticateUser(ctx context.Context, user dto.AuthRequest) (entity.Session, error) {
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
