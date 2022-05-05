package impl

import (
	"context"
	"database/sql"
	"log"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type UserAuthRepositoryImpl struct {
	DB *sql.DB
}

var (
	INSERT_USER = "INSERT INTO users(username, password) VALUE (?, ?)"
	FIND_USER   = "SELECT uid, username, password FROM users WHERE username = ?"
)

func NewUserAuthRepository(DB *sql.DB) UserAuthRepositoryImpl {
	return UserAuthRepositoryImpl{
		DB: DB,
	}
}

func (u UserAuthRepositoryImpl) NewUser(ctx context.Context, user entity.User) error {
	prpd, err := u.DB.PrepareContext(ctx, INSERT_USER)
	if err != nil {
		log.Println("[ERROR] NewUser -> error :", err)
		return err
	}

	result, err := prpd.ExecContext(ctx, user.Username, user.Password)
	if err != nil {
		log.Println("[ERROR] NewUser -> error on executing query :", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] NewUser -> error on getting rows affected :", err)
		return err
	}
	if rows != 1 {
		log.Println("[ERROR] NewUser -> error on inserting row :", err)
		return err
	}

	return nil
}

func (u UserAuthRepositoryImpl) GetUser(ctx context.Context, username string) (entity.User, error) {
	prpd, err := u.DB.PrepareContext(ctx, FIND_USER)
	if err != nil {
		log.Println("[ERROR] GetUser -> error :", err)
		return entity.User{}, err
	}

	rows, err := prpd.Query(username)
	if err != nil {
		log.Println("[ERROR] GetUser -> error on executing query :", err)
		return entity.User{}, err
	}

	var user entity.User
	if rows.Next() {
		err = rows.Scan(&user.UserID, &user.Username, &user.Password)
		if err != nil {
			log.Println("[ERROR] GetUser -> error scanning row :", err)
			return entity.User{}, err
		}

		return user, nil
	}

	return entity.User{}, err
}
