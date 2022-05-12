package impl

import (
	"context"
	"database/sql"
	"log"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type CommentRepositoryImpl struct {
	db *sql.DB
}

var (
	INSERT_COMMENT = "INSERT INTO comment(post_id, uid, comment_body) VALUE (?, ?, ?)"
)

func NewCommentRepository(db *sql.DB) CommentRepositoryImpl {
	return CommentRepositoryImpl{
		db: db,
	}
}

func (c CommentRepositoryImpl) AddComment(ctx context.Context, comment entity.Comment) error {
	result, err := c.db.ExecContext(ctx, INSERT_COMMENT, comment.PostID, comment.UserID, comment.CommentBody)
	if err != nil {
		log.Println("[ERROR] AddComment -> error inserting row :", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] AddComment -> error on getting rows affected :", err)
		return err
	}
	if rowsAffected != 1 {
		log.Println("[ERROR] AddComment -> error on updating row :", err)
		return err
	}

	return nil
}
