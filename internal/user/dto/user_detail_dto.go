package dto

import (
	"time"

	"github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type MiniUserDetailResponse struct {
	UserID    int64  `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   string `json:"picture"`
}

type UserDetailResponse struct {
	UserID    int64          `json:"uid"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Picture   string         `json:"picture"`
	Phone     string         `json:"phone"`
	AboutMe   dto.NullString `json:"about_me"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type UserDetailRequest struct {
	UserID    int64  `json:"uid"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   string `json:"picture"`
	Phone     string `json:"phone"`
	AboutMe   string `json:"about_me"`
}

func NewMiniUserDetailDTO(user entity.MiniUserDetail) MiniUserDetailResponse {
	return MiniUserDetailResponse{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Picture:   user.Picture,
	}
}

func NewUserDetailDTO(user entity.UserDetail) UserDetailResponse {
	return UserDetailResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Picture:   user.Picture,
		Phone:     user.Phone,
		AboutMe:   user.AboutMe,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewUserDetailDAO(user UserDetailRequest) entity.UserDetail {
	return entity.UserDetail{
		UserID:    user.UserID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Picture:   user.Picture,
		Phone:     user.Phone,
		AboutMe: dto.NullString{
			String: user.AboutMe,
			Valid:  true,
		},
	}
}
