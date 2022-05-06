package impl

import (
	"context"
	"database/sql"
	"log"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type PostRepositoryImpl struct {
	DB *sql.DB
}

var (
	INSERT_POST = "INSERT INTO users(username, password) VALUE (?, ?)"
	SELECT_POST = "SELECT uid, username, password FROM users WHERE username = ?"
)

func NewPostRepository(DB *sql.DB) PostRepositoryImpl {
	return PostRepositoryImpl{
		DB: DB,
	}
}

func (p PostRepositoryImpl) NewPost(ctx context.Context, user entity.User) error {
	prpd, err := p.DB.PrepareContext(ctx, INSERT_POST)
	if err != nil {
		log.Println("[ERROR] NewPost -> error :", err)
		return err
	}

	result, err := prpd.ExecContext(ctx, user.Username, user.Password)
	if err != nil {
		log.Println("[ERROR] NewPost -> error on executing query :", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] NewPost -> error on getting rows affected :", err)
		return err
	}
	if rows != 1 {
		log.Println("[ERROR] NewPost -> error on inserting row :", err)
		return err
	}

	return nil
}

func (p PostRepositoryImpl) GetPost(ctx context.Context, username string) (entity.User, error) {
	prpd, err := p.DB.PrepareContext(ctx, SELECT_POST)
	if err != nil {
		log.Println("[ERROR] GetPost -> error :", err)
		return entity.User{}, err
	}

	rows, err := prpd.Query(username)
	if err != nil {
		log.Println("[ERROR] GetPost -> error on executing query :", err)
		return entity.User{}, err
	}

	var user entity.User
	if rows.Next() {
		err = rows.Scan(&user.UserID, &user.Username, &user.Password)
		if err != nil {
			log.Println("[ERROR] GetPost -> error scanning row :", err)
			return entity.User{}, err
		}

		return user, nil
	}

	return entity.User{}, err
}
