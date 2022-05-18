package impl

import (
	"context"
	"database/sql"
	"log"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type SessionRepositoryImpl struct {
	db *sql.DB
}

var (
	INSERT_SESSION = "INSERT INTO sessions(token, uid, expireAt) VALUE (?, ?, ?)"
	FIND_SESSION   = "SELECT s.token, s.uid, u.username, s.expireAt FROM sessions s JOIN users u ON s.uid = u.uid WHERE token = ?"
	DELETE_SESSION = "DELETE FROM sessions WHERE token = ?"
)

func NewSessionRepository(db *sql.DB) SessionRepositoryImpl {
	return SessionRepositoryImpl{
		db: db,
	}
}

func (u SessionRepositoryImpl) NewSession(ctx context.Context, user entity.Session) error {
	result, err := u.db.ExecContext(ctx, INSERT_SESSION, user.Token, user.UID, user.ExpireAt)
	if err != nil {
		log.Println("[ERROR] NewSession -> error on executing query :", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] NewSession -> error on getting rows affected :", err)
		return err
	}
	if rows != 1 {
		log.Println("[ERROR] NewSession -> error on inserting row :", err)
		return err
	}

	return nil
}

func (u SessionRepositoryImpl) GetSession(ctx context.Context, token string) (entity.SessionDetail, error) {
	rows, err := u.db.QueryContext(ctx, FIND_SESSION, token)
	if err != nil {
		log.Println("[ERROR] GetSession -> error on executing query :", err)
		return entity.SessionDetail{}, err
	}

	var user entity.SessionDetail
	if rows.Next() {
		err = rows.Scan(&user.Token, &user.UID, &user.Username, &user.ExpireAt)
		if err != nil {
			log.Println("[ERROR] GetSession -> error scanning row :", err)
			return entity.SessionDetail{}, err
		}

		return user, nil
	}

	return entity.SessionDetail{}, err
}

func (u SessionRepositoryImpl) DeleteSession(ctx context.Context, token string) error {
	result, err := u.db.ExecContext(ctx, DELETE_SESSION, token)
	if err != nil {
		log.Println("[ERROR] DeleteSession -> error on executing query :", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] DeleteSession -> error on getting rows affected :", err)
		return err
	}
	if rows != 1 {
		log.Println("[ERROR] DeleteSession -> error on deleting row :", err)
		return err
	}

	return nil
}
