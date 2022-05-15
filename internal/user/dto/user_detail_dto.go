package dto

import (
	"time"

	"github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type MiniUserDetail struct {
	UserID    int64  `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   string `json:"picture"`
}

type UserDetail struct {
	UserID    int64          `json:"uid"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Picture   string         `json:"picture"`
	Phone     string         `json:"phone"`
	AboutMe   dto.NullString `json:"about_me"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func NewMiniUserDetailDTO(user entity.MiniUserDetail) MiniUserDetail {
	return MiniUserDetail{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Picture:   user.Picture,
	}
}

func NewUserDetailDTO(user entity.UserDetail) UserDetail {
	return UserDetail{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Picture:   user.Picture,
		Phone:     user.Phone,
		AboutMe:   user.AboutMe,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
