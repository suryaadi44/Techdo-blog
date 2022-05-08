package dto

import "github.com/suryaadi44/Techdo-blog/pkg/entity"

type MiniUserDetail struct {
	UserID    int64  `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   string `json:"picture"`
}

func NewMiniUserDetailDTO(user entity.UserDetail) MiniUserDetail {
	return MiniUserDetail{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Picture:   user.Picture,
	}
}
