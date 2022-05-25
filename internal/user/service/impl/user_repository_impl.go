package impl

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

var (
	COUNT_TOTAL_USER_POST    = "SELECT COUNT(*) FROM blog_posts WHERE author_id = ?"
	COUNT_TOTAL_USER_COMMENT = "SELECT COUNT(*) FROM comment WHERE uid = ?"

	INSERT_USER_DETAIL = "INSERT INTO user_details(uid, email, first_name, last_name, picture, phone, about_me) VALUE (?, ?, ?, ?, ?, ?, ?)"

	UPDATE_USER_DETAIL  = "UPDATE user_details SET first_name = ?, last_name = ?, phone = ?, about_me = ? WHERE uid = ?"
	UPDATE_USER_PICTURE = "UPDATE user_details SET picture = ? WHERE uid = ?"

	SELECT_USER_DETAIL          = "SELECT d.uid, u.username, d.email, d.first_name, d.last_name, d.picture, d.phone, d.about_me, d.created_at, d.updated_at FROM user_details d JOIN users u ON d.uid = u.uid WHERE d.uid = ?"
	SELECT_USER_MINI_DETAIL     = "SELECT d.uid, u.username, d.first_name, d.last_name, d.picture FROM user_details d JOIN users u ON d.uid = u.uid WHERE d.uid = ?"
	SELECT_USER_PICTURE_PROFILE = "SELECT picture FROM user_details WHERE uid = ?"

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
		log.Println("[ERROR] UpdateUserDetail -> error on updating row :", err)
		return err
	}

	return nil
}
func (u UserRepositoryImpl) UpdateUserPicture(ctx context.Context, url string, id int64) error {
	result, err := u.db.ExecContext(ctx, UPDATE_USER_PICTURE, url, id)
	if err != nil {
		log.Println("[ERROR] UpdateUserPicture -> error on executing query :", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] UpdateUserPicture -> error on getting rows affected :", err)
		return err
	}
	if rows != 1 {
		log.Println("[ERROR] UpdateUserPicture -> error on updating row :", err)
		return err
	}

	return nil
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

	return user, errors.New("can't get user detail")
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

	return user, errors.New("can't get user mini detail")
}

func (u UserRepositoryImpl) GetUserPictureID(ctx context.Context, id int64) (string, error) {
	rows, err := u.db.QueryContext(ctx, SELECT_USER_PICTURE_PROFILE, id)
	if err != nil {
		log.Println("[ERROR] GetUserPictureID -> error on executing query :", err)
		return "", err
	}

	if rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			log.Println("[ERROR] GetUserPictureID -> error scanning row :", err)
			return "", err
		}

		return id, nil
	}

	return "", errors.New("can't get user picture")
}

func (u UserRepositoryImpl) GetUserTotalPostCount(ctx context.Context, id int64) (int64, error) {
	var count int64

	rows, err := u.db.QueryContext(ctx, COUNT_TOTAL_USER_POST, id)
	if err != nil {
		log.Println("[ERROR] GetUserTotalPostCount -> error on executing query :", err)
		return 0, err
	}

	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println("[ERROR] GetUserTotalPostCount -> error scanning row :", err)
			return 0, err
		}

		return count, nil
	}

	return 0, errors.New("can't get count of total user post")
}
func (u UserRepositoryImpl) GetUserTotalCommentCount(ctx context.Context, id int64) (int64, error) {
	var count int64

	rows, err := u.db.QueryContext(ctx, COUNT_TOTAL_USER_COMMENT, id)
	if err != nil {
		log.Println("[ERROR] GetUserTotalCommentCount -> error on executing query :", err)
		return 0, err
	}

	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println("[ERROR] GetUserTotalCommentCount -> error scanning row :", err)
			return 0, err
		}

		return count, nil
	}

	return 0, errors.New("can't get count of total user comment")
}
