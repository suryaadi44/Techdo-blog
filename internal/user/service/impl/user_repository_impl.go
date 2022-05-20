package impl

import (
	"context"
	"database/sql"
	"log"

	"github.com/codedius/imagekit-go"
	"github.com/suryaadi44/Techdo-blog/pkg/entity"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

var (
	INSERT_USER_DETAIL = "INSERT INTO user_details(uid, email, first_name, last_name, picture, phone, about_me) VALUE (?, ?, ?, ?, ?, ?, ?)"

	UPDATE_USER_DETAIL = "UPDATE user_details SET first_name = ?, last_name = ?, phone = ?, about_me = ? WHERE uid = ?"

	SELECT_USER_DETAIL      = "SELECT d.uid, u.username, d.email, d.first_name, d.last_name, d.picture, d.phone, d.about_me, d.created_at, d.updated_at FROM user_details d JOIN users u ON d.uid = u.uid WHERE d.uid = ?"
	SELECT_USER_MINI_DETAIL = "SELECT d.uid, u.username, d.first_name, d.last_name, d.picture FROM user_details d JOIN users u ON d.uid = u.uid WHERE d.uid = ?"

	DELETE_USER = "DELETE FROM users WHERE uid = ?"
)

func NewUserRepository(db *sql.DB) UserRepositoryImpl {
	return UserRepositoryImpl{
		db: db,
	}
}

func (u UserRepositoryImpl) DeleteUser(ctx context.Context, id int64) error {
	result, err := u.db.ExecContext(ctx, DELETE_USER, id)
	if err != nil {
		log.Println("[ERROR] DeletUser -> error on executing query :", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] DeletUser -> error on getting rows affected :", err)
		return err
	}
	if rows != 1 {
		log.Println("[ERROR] DeletUser -> error on deleting row :", err)
		return err
	}

	return nil
}

func (u UserRepositoryImpl) UpdateUserDetail(ctx context.Context, user entity.UserDetail) error {
	result, err := u.db.ExecContext(ctx, UPDATE_USER_DETAIL, user.FirstName, user.LastName, user.Phone, user.AboutMe.String, user.UserID)
	if err != nil {
		log.Println("[ERROR] UpdateUserDetail -> error on executing query :", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] UpdateUserDetail -> error on getting rows affected :", err)
		return err
	}
	if rows != 1 {
		log.Println("[ERROR] UpdateUserDetail -> error on Updating row :", err)
		return err
	}

	return nil
}

func (u UserRepositoryImpl) UploadImage(ctx context.Context, filename string, image interface{}, id string) (*imagekit.UploadResponse, error) {
	folderID := "/user"
	return utils.UploadImage(ctx, id, image, folderID)
}

func (u UserRepositoryImpl) GetUserDetail(ctx context.Context, id int64) (entity.UserDetail, error) {
	var user entity.UserDetail
	var username string

	rows, err := u.db.QueryContext(ctx, SELECT_USER_DETAIL, id)
	if err != nil {
		log.Println("[ERROR] GetUserDetail -> error on executing query :", err)
		return user, err
	}

	if rows.Next() {
		err = rows.Scan(&user.UserID, &username, &user.Email, &user.FirstName, &user.LastName, &user.Picture, &user.Phone, &user.AboutMe, &user.CreatedAt, &user.UpdatedAt)
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
		err = rows.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Picture)
		if err != nil {
			log.Println("[ERROR] GetUserMiniDetail -> error scanning row :", err)
			return user, err
		}

		return user, nil
	}

	return user, err
}
