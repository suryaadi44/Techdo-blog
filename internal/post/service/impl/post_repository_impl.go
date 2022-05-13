package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/suryaadi44/Techdo-blog/pkg/entity"
)

type PostRepositoryImpl struct {
	db *sql.DB
}

var (
	INSERT_BLANK_POST     = "INSERT INTO blog_posts() VALUE ()"
	INSERT_CATEGORY_ASSOC = "INSERT INTO category_associations(post_id, category_id) VALUE (?, ?)"
	INSERT_COMMENT        = "INSERT INTO comment(post_id, uid, comment_body) VALUE (?, ?, ?) ORDER BY b.created_at DESC LIMIT ?, ?"

	UPDATE_POST = "UPDATE blog_posts SET author_id = ?, banner = ?, title = ?, body = ? WHERE post_id = ?"

	DELETE_POST = "DELETE FROM blog_posts WHERE post_id = ?"

	SELECT_POST = "SELECT b.post_id, b.author_id, b.banner, b.title, b.body, b.created_at, b.updated_at, u.uid, u.email, u.first_name, u.last_name, u.picture, u.phone, u.about_me, u.created_at, u.updated_at FROM blog_posts b JOIN user_details u ON b.author_id = u.uid WHERE b.post_id = ?"

	SELECT_POST_AUTHOR       = "SELECT author_id FROM blog_posts WHERE post_id = ?"
	SELECT_ID_OF_LAST_INSERT = "SELECT LAST_INSERT_ID() as uid"
	SELECT_LIST_OF_POST      = "SELECT b.post_id, b.banner, b.title, b.body, b.created_at, b.updated_at, CONCAT(u.first_name, ' ', u.last_name) AS author FROM blog_posts b JOIN user_details u ON b.author_id = u.uid"
	SELECT_FULL_TEXT_POST    = "SELECT b.post_id, b.banner, b.title, b.body, b.created_at, b.updated_at, CONCAT(u.first_name, ' ', u.last_name) AS author FROM blog_posts b JOIN user_details u ON b.author_id = u.uid WHERE MATCH(b.title) AGAINST(? IN NATURAL LANGUAGE MODE)"
	SELECT_CATEGORY_OF_POST  = "SELECT c.category_id, c.category_name FROM categories c JOIN category_associations a ON c.category_id = a.category_id WHERE a.post_id = ?"
	SELECT_CATEGORY          = "SELECT category_id, category_name FROM categories"
	SELECT_COMMENTS          = "SELECT comment_id, uid, comment_body, created_at, updated_at WHERE post_id = ?"
)

func NewPostRepository(db *sql.DB) PostRepositoryImpl {
	return PostRepositoryImpl{
		db: db,
	}
}

func (p PostRepositoryImpl) ReserveID(ctx context.Context) (int64, error) {
	res, err := p.db.ExecContext(ctx, INSERT_BLANK_POST)
	if err != nil {
		log.Println("[ERROR] ReserveID -> error inserting blank row :", err)
		return -1, err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		log.Println("[ERROR] ReserveID -> error getting row :", err)
		return -1, err
	}

	return lid, nil
}

func (p PostRepositoryImpl) UpdatePost(ctx context.Context, post entity.BlogPost) error {
	result, err := p.db.ExecContext(ctx, UPDATE_POST, post.AuthorID, post.Banner, post.Title, post.Body, post.PostID)
	if err != nil {
		log.Println("[ERROR] UpdatePost -> error on executing query :", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] UpdatePost -> error on getting rows affected :", err)
		return err
	}
	if rowsAffected != 1 {
		log.Println("[ERROR] UpdatePost -> error on updating row :", err)
		return err
	}

	return nil
}

func (p PostRepositoryImpl) DeletePost(ctx context.Context, id int64) error {
	result, err := p.db.ExecContext(ctx, DELETE_POST, id)
	if err != nil {
		log.Println("[ERROR] DeletePost -> error on executing query :", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] DeletePost -> error on getting rows affected :", err)
		return err
	}
	if rowsAffected != 1 {
		log.Println("[ERROR] DeletePost -> error on deleting row :", err)
		return err
	}

	return nil
}

func (p PostRepositoryImpl) GetFullPost(ctx context.Context, id int64) (entity.BlogPost, entity.UserDetail, error) {
	var post entity.BlogPost
	var author entity.UserDetail

	rows, err := p.db.QueryContext(ctx, SELECT_POST, id)
	if err != nil {
		log.Println("[ERROR] GetFullPost -> error on executing query :", err)
		return post, author, err
	}

	if rows.Next() {
		err = rows.Scan(&post.PostID, &post.AuthorID, &post.Banner, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt, &author.UserID, &author.Email, &author.FirstName, &author.LastName, &author.Picture, &author.Phone, &author.AboutMe, &author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			log.Println("[ERROR] GetFullPost -> error scanning row :", err)
			return post, author, err
		}

		return post, author, nil
	}

	return post, author, fmt.Errorf("No post with id %d", id)
}

func (p PostRepositoryImpl) GetBriefsBlogPostData(ctx context.Context, offset int64, limit int64) (entity.BriefsBlogPost, error) {
	var postList entity.BriefsBlogPost

	query := SELECT_LIST_OF_POST + " ORDER BY b.created_at DESC LIMIT ?, ? "
	rows, err := p.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		log.Println("[ERROR] GetBriefsBlogPostData -> error on executing query :", err)
		return postList, err
	}

	for rows.Next() {
		var post entity.BriefBlogPost
		err = rows.Scan(&post.PostID, &post.Banner, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			log.Println("[ERROR] GetBriefsBlogPostData -> error scanning row :", err)
			return postList, err
		}

		postList = append(postList, &post)
	}

	return postList, nil
}

func (p PostRepositoryImpl) GetBriefsBlogPostFromSearch(ctx context.Context, q string, offset int64, limit int64, dateStart *time.Time, dateEnd *time.Time) (entity.BriefsBlogPost, error) {
	var postList entity.BriefsBlogPost
	var query string
	var args []interface{}

	args = append(args, q)
	query = SELECT_FULL_TEXT_POST

	if dateStart != nil && dateEnd != nil {
		query = query + " AND b.created_at BETWEEN ? AND ?"
		args = append(args, dateStart.Format("2006-01-02"), dateEnd.Format("2006-01-02"))
	} else if dateStart != nil {
		query = query + " AND b.created_at > ?"
		args = append(args, dateStart.Format("2006-01-02"))
	} else if dateEnd != nil {
		query = query + " AND b.created_at < ?"
		args = append(args, dateEnd.Format("2006-01-02"))
	}

	query = query + " ORDER BY b.created_at DESC LIMIT ?, ?"
	args = append(args, offset, limit)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println("[ERROR] GetBriefsBlogPostFromSearch -> error on executing query :", err)
		return postList, err
	}

	for rows.Next() {
		var post entity.BriefBlogPost
		err = rows.Scan(&post.PostID, &post.Banner, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			log.Println("[ERROR] GetBriefsBlogPostFromSearch -> error scanning row :", err)
			return postList, err
		}

		postList = append(postList, &post)
	}

	return postList, nil
}

func (p PostRepositoryImpl) GetPostAuthorId(ctx context.Context, postID int64) (int64, error) {
	rows, err := p.db.QueryContext(ctx, SELECT_POST_AUTHOR, postID)
	if err != nil {
		log.Println("[ERROR] GetPostAuthorId -> error on executing query :", err)
		return -1, err
	}

	var authorID int64
	if rows.Next() {
		err = rows.Scan(&authorID)
		if err != nil {
			log.Println("[ERROR] GetPostAuthorId -> error scanning row :", err)
			return -1, err
		}

		return authorID, nil
	}

	return -1, fmt.Errorf("No post with id %d", postID)
}

func (p PostRepositoryImpl) GetCategoriesFromID(ctx context.Context, id int64) (entity.Categories, error) {
	var categories entity.Categories

	rows, err := p.db.QueryContext(ctx, SELECT_CATEGORY_OF_POST, id)
	if err != nil {
		log.Println("[ERROR] GetCategoriesFromID -> error on executing query :", err)
		return categories, err
	}

	found := false
	for rows.Next() {
		found = true
		var postCategory entity.Category

		err = rows.Scan(&postCategory.CategoryID, &postCategory.CategoryName)
		if err != nil {
			log.Println("[ERROR] GetCategoriesFromID -> error scanning row :", err)
			return categories, err
		}
		categories = append(categories, &postCategory)
	}

	if !found {
		return categories, errors.New("Categories not found")
	}

	return categories, nil
}

func (p PostRepositoryImpl) GetCategoryList(ctx context.Context) (entity.Categories, error) {
	var categories entity.Categories

	rows, err := p.db.QueryContext(ctx, SELECT_CATEGORY)
	if err != nil {
		log.Println("[ERROR] GetCategoryList -> error on executing query :", err)
		return categories, err
	}

	for rows.Next() {
		var postCategory entity.Category

		err = rows.Scan(&postCategory.CategoryID, &postCategory.CategoryName)
		if err != nil {
			log.Println("[ERROR] GetCategoryList -> error scanning row :", err)
			return categories, err
		}
		categories = append(categories, &postCategory)

	}
	return categories, nil
}

func (p PostRepositoryImpl) AddPostCategoryAssoc(ctx context.Context, posdtId int64, categoryId int64) error {
	result, err := p.db.ExecContext(ctx, INSERT_CATEGORY_ASSOC, posdtId, categoryId)
	if err != nil {
		log.Println("[ERROR] AddPostCategoryAssoc -> error inserting row :", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] AddPostCategoryAssoc -> error on getting rows affected :", err)
		return err
	}
	if rowsAffected != 1 {
		log.Println("[ERROR] AddPostCategoryAssoc -> error on updating row :", err)
		return err
	}

	return nil
}

func (p PostRepositoryImpl) AddComment(ctx context.Context, comment entity.Comment) error {
	result, err := p.db.ExecContext(ctx, INSERT_COMMENT, comment.PostID, comment.UserID, comment.CommentBody)
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

func (p PostRepositoryImpl) GetPostComments(ctx context.Context, id int64, offset int64, limit int64) (entity.Comments, entity.MiniUsersDetail, error) {
	var comments entity.Comments
	var users entity.MiniUsersDetail

	rows, err := p.db.QueryContext(ctx, SELECT_COMMENTS, id, offset, limit)
	if err != nil {
		log.Println("[ERROR] GetPostComments -> error on executing query :", err)
		return comments, users, err
	}

	for rows.Next() {
		var comment entity.Comment
		err = rows.Scan(&comment.CommentID, &comment.UserID, &comment.CommentBody, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			log.Println("[ERROR] GetPostComments -> error scanning row :", err)
			return comments, users, err
		}

		comments = append(comments, &comment)
	}

	return comments, users, nil
}
