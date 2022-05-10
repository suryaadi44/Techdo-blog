package impl

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type UserAuthRepositoryImpl struct {
	DB *sql.DB
}

var (
	INSERT_USER        = "INSERT INTO users(username, password) VALUE (?, ?)"
	INSERT_USER_DETAIL = "INSERT INTO user_details(uid, email, first_name, last_name, picture) VALUE (?, ?, ?, ?, ?)"
	FIND_USER          = "SELECT uid, username, password FROM users WHERE username = ?"
)

func NewUserAuthRepository(DB *sql.DB) UserAuthRepositoryImpl {
	return UserAuthRepositoryImpl{
		DB: DB,
	}
}

func (u UserAuthRepositoryImpl) NewUser(ctx context.Context, user entity.User, userDetail entity.UserDetail) error {
	result, err := u.DB.ExecContext(ctx, INSERT_USER, user.Username, user.Password)
	if err != nil {
		log.Println("[ERROR] NewUser -> error on executing insert user query :", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] NewUser -> error on getting insert user rows affected :", err)
		return err
	}

	if rows != 1 {
		log.Println("[ERROR] NewUser -> error on inserting insert user row :", err)
		return errors.New("Cant creare new user")
	}

	lid, err := result.LastInsertId()
	if err != nil {
		log.Println("[ERROR] NewUser -> error on getting uid :", err)
		return err
	}

	result, err = u.DB.ExecContext(ctx, INSERT_USER_DETAIL, lid, userDetail.Email, userDetail.FirstName, userDetail.LastName, userDetail.Picture)
	if err != nil {
		log.Println("[ERROR] NewUser -> error on executing insert user details query :", err)
		return err
	}

	rows, err = result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] NewUser -> error on getting insert user details rows affected :", err)
		return err
	}

	if rows != 1 {
		log.Println("[ERROR] NewUser -> error on inserting insert user details row  :", err)
		return errors.New("Cant creare new user")
	}

	return nil
}

func (u UserAuthRepositoryImpl) GetUser(ctx context.Context, username string) (entity.User, error) {
	var user entity.User

	prpd, err := u.DB.PrepareContext(ctx, FIND_USER)
	if err != nil {
		log.Println("[ERROR] GetUser -> error :", err)
		return user, err
	}

	rows, err := prpd.Query(username)
	if err != nil {
		log.Println("[ERROR] GetUser -> error on executing query :", err)
		return user, err
	}

	if rows.Next() {
		err = rows.Scan(&user.UserID, &user.Username, &user.Password)
		if err != nil {
			log.Println("[ERROR] GetUser -> error scanning row :", err)
			return entity.User{}, err
		}

		return user, nil
	}

	return user, errors.New("Inccorect Username or Password")
}
