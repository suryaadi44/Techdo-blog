package impl

import (
	"context"
	"database/sql"
	"log"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

var (
	INSERT_USER_DETAIL = "INSERT INTO user_details(uid, email, first_name, last_name, picture, phone, about_me) VALUE (?, ?, ?, ?, ?, ?, ?)"

	SELECT_USER_DETAIL      = "SELECT uid, email, first_name, last_name, picture, phone, about_me, created_at, updated_at FROM user_details WHERE uid = ?"
	SELECT_USER_MINI_DETAIL = "SELECT uid, first_name, last_name, picture FROM user_details WHERE uid = ?"
)

func NewUserRepository(db *sql.DB) UserRepositoryImpl {
	return UserRepositoryImpl{
		db: db,
	}
}

// func (u UserRepositoryImpl) NewUserDetail(ctx context.Context, id int64) error {
// 	prpd, err := u.DB.PrepareContext(ctx, INSERT_USER_DETAIL)
// 	if err != nil {
// 		log.Println("[ERROR] NewUserDetail -> error :", err)
// 		return err
// 	}

// 	result, err := prpd.ExecContext(ctx)
// 	if err != nil {
// 		log.Println("[ERROR] NewUserDetail -> error on executing query :", err)
// 		return err
// 	}

// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		log.Println("[ERROR] NewUserDetail -> error on getting rows affected :", err)
// 		return err
// 	}
// 	if rows != 1 {
// 		log.Println("[ERROR] NewUserDetail -> error on inserting row :", err)
// 		return err
// 	}

// 	return nil
// }

func (u UserRepositoryImpl) GetUserDetail(ctx context.Context, id int64) (entity.UserDetail, error) {
	var user entity.UserDetail

	rows, err := u.db.QueryContext(ctx, SELECT_USER_DETAIL, id)
	if err != nil {
		log.Println("[ERROR] GetUserDetail -> error on executing query :", err)
		return user, err
	}

	if rows.Next() {
		err = rows.Scan(&user.UserID, &user.Email, &user.FirstName, &user.LastName, &user.Picture, &user.Phone, &user.AboutMe, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println("[ERROR] GetUserDetail -> error scanning row :", err)
			return user, err
		}

		return user, nil
	}

	return user, err
}

func (u UserRepositoryImpl) GetUserMiniDetail(ctx context.Context, id int64) (entity.MiniUserDetail, error) {
	var user entity.MiniUserDetail

	rows, err := u.db.QueryContext(ctx, SELECT_USER_MINI_DETAIL, id)
	if err != nil {
		log.Println("[ERROR] GetUserMiniDetail -> error on executing query :", err)
		return user, err
	}

	if rows.Next() {
		err = rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Picture)
		if err != nil {
			log.Println("[ERROR] GetUserMiniDetail -> error scanning row :", err)
			return user, err
		}

		return user, nil
	}

	return user, err
}
