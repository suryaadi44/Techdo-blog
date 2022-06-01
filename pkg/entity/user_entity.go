package entity

import (
	"time"

	"github.com/suryaadi44/Techdo-blog/pkg/dto"
)

type User struct {
	UserID   int64  `db:"uid"`
	Username string `db:"username"`
	Password string `db:"password"`
	Type     int64  `db:"type"`
}

type UserDetail struct {
	UserID    int64          `db:"uid"`
	Email     string         `db:"email"`
	FirstName string         `db:"first_name"`
	LastName  string         `db:"last_name"`
	Picture   string         `db:"picture"`
	Phone     string         `db:"phone"`
	AboutMe   dto.NullString `db:"about_me"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type MiniUserDetail struct {
	UserID    int64     `db:"uid"`
	Username  string    `db:"username"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Picture   string    `db:"picture"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Session struct {
	Token    string    `db:"token"`
	UID      int64     `db:"uid"`
	ExpireAt time.Time `db:"expireAt"`
}

type SessionDetail struct {
	Token    string    `db:"token"`
	UID      int64     `db:"uid"`
	Username string    `db:"username"`
	ExpireAt time.Time `db:"expireAt"`
}

type MiniUsersDetail []*MiniUserDetail
