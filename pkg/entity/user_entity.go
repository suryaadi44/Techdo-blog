package entity

import "time"

type User struct {
	UserID   int64  `db:"uid"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserDetail struct {
	UserID    int64     `db:"uid"`
	Email     string    `db:"email"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Picture   string    `db:"picture"`
	Phone     string    `db:"phone"`
	AboutMe   string    `db:"about_me"`
	CreatedAt time.Time `db:"created_at"`
}

type Session struct {
	Token    string    `db:"token"`
	UID      int64     `db:"uid"`
	ExpireAt time.Time `db:"expireAt"`
}
