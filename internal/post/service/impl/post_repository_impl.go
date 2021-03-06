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
	ADD_VIEW = `UPDATE blog_posts 
				SET view_count = view_count + 1 
				WHERE post_id = ?`

	COUNT_LIST_OF_POST               = `SELECT COUNT(*) 
										FROM blog_posts`
	COUNT_LIST_OF_POST_IN_CATEGOIES  = `SELECT COUNT(*) 
										FROM blog_posts b 
										JOIN category_associations a 
										ON  a.post_id =  b.post_id 
										JOIN categories c 
										ON c.category_id = a.category_id 
										WHERE c.category_name = ?`
	COUNT_SEARCH_RESULT              = `SELECT COUNT(*) 
										FROM blog_posts b`
	COUNT_TOTAL_USER_POST            = `SELECT COUNT(*) 
										FROM blog_posts 
										WHERE author_id = ?`
	COUNT_TOTAL_USER_COMMENT         = `SELECT COUNT(*) 
										FROM comment 
										WHERE uid = ?`

	INSERT_BLANK_POST     = `INSERT INTO blog_posts() 
							VALUE ()`
	INSERT_CATEGORY_ASSOC = `INSERT INTO category_associations(post_id, category_id) 
							VALUE (?, ?)`
	INSERT_COMMENT        = `INSERT INTO comment(post_id, uid, comment_body) 
							VALUE (?, ?, ?)`

	UPDATE_POST           = `UPDATE blog_posts 
							SET author_id = ?, banner = ?, title = ?, body = ? 
							WHERE post_id = ?`
	UPDATE_CATEGORY_ASSOC = `UPDATE category_associations 
							SET category_id = ? 
							WHERE post_id = ?`

	DELETE_POST           = `DELETE FROM blog_posts 
							WHERE post_id = ?`
	DELETE_COMMENT        = `DELETE FROM comment 
							WHERE comment_id = ?`
	DELETE_CATEGORY_ASSOC = `DELETE FROM category_associations 
							WHERE post_id = ?`


	SELECT_RAW_POST                          = `SELECT post_id, author_id, banner, title, body, created_at, updated_at 
												FROM blog_posts 
												WHERE post_id = ?`
	SELECT_POST                              = `SELECT b.post_id, b.author_id, b.banner, b.title, b.body, b.view_count, b.comment_count, b.created_at, b.updated_at, u.uid, u.email, u.first_name, u.last_name, u.picture, u.phone, u.about_me, u.created_at, u.updated_at 
												FROM blog_posts b 
												JOIN user_details u 
												ON b.author_id = u.uid 
												WHERE b.post_id = ?`
	SELECT_POST_AUTHOR                       = `SELECT author_id 
												FROM blog_posts 
												WHERE post_id = ?`
	SELECT_COMMENT_AUTHOR                    = `SELECT uid 
												FROM comment 
												WHERE comment_id = ?`
	SELECT_ID_OF_LAST_INSERT                 = "SELECT LAST_INSERT_ID() as uid"
	SELECT_LIST_OF_POST                      = `SELECT b.post_id, b.banner, b.title, b.body, b.view_count, b.comment_count, b.created_at, b.updated_at, CONCAT(u.first_name, ' ', u.last_name) AS author 
												FROM blog_posts b 
												JOIN user_details u 
												ON b.author_id = u.uid 
												ORDER BY b.created_at DESC LIMIT ?, ? `
	SELECT_LIST_OF_POST_BY_USER              = `SELECT b.post_id, b.title, b.created_at, c.category_name 
												FROM blog_posts b LEFT JOIN category_associations a 
												ON a.post_id = b.post_id 
												LEFT JOIN categories c 
												ON a.category_id = c.category_id 
												WHERE b.author_id = ? OR b.author_id IS NULL 
												UNION 
												SELECT b.post_id, b.title, b.created_at, c.category_name 
												FROM blog_posts b RIGHT JOIN category_associations a 
												ON a.post_id = b.post_id 
												RIGHT JOIN categories c 
												ON a.category_id = c.category_id 
												WHERE b.author_id = ? OR b.author_id IS NULL 
												ORDER BY created_at DESC`
	SELECT_EACH_CATEGORY_USER_POST_STATISTIC = `SELECT c.category_name, COUNT(b.post_id) AS TotalPost, SUM(b.view_count) AS TotalView 
												FROM blog_posts b 
												LEFT JOIN category_associations a 
												ON a.post_id = b.post_id 
												LEFT JOIN categories c 
												ON a.category_id = c.category_id 
												WHERE b.author_id = ? 
												GROUP BY c.category_name 
													HAVING TotalView > 0 
												ORDER BY TotalView DESC`
	SELECT_LISF_OF_POST_IN_CATEGORY          = `SELECT b.post_id, b.banner, b.title, b.body, b.view_count, b.comment_count, b.created_at, b.updated_at, CONCAT(u.first_name, ' ', u.last_name) AS author 
												FROM blog_posts b JOIN user_details u 
												ON b.author_id = u.uid 
												JOIN category_associations a 
												ON  a.post_id =  b.post_id 
												JOIN categories c 
												ON c.category_id = a.category_id 
												WHERE c.category_name = ? 
												ORDER BY b.created_at DESC LIMIT ?, ?`
	SELECT_FULL_TEXT_POST                    = `SELECT b.post_id, b.banner, b.title, b.body, b.view_count, b.comment_count, b.created_at, b.updated_at, CONCAT(u.first_name, ' ', u.last_name) AS author 
												FROM blog_posts b 
												JOIN user_details u 
												ON b.author_id = u.uid`
	SELECT_CATEGORY_OF_POST                  = `SELECT c.category_id, c.category_name 
												FROM categories c 
												JOIN category_associations a 
												ON c.category_id = a.category_id 
												WHERE a.post_id = ?`

	SELECT_CATEGORY                        	 = `SELECT category_id, category_name 
												FROM categories`
	SELECT_COMMENTS                        	 = `SELECT c.comment_id, c.uid, c.comment_body, c.created_at, c.updated_at, u.uid, u.first_name, u.last_name, u.picture 
												FROM comment c 
												JOIN user_details u 
												ON c.uid= u.uid 
												WHERE c.post_id = ? 
												ORDER BY c.created_at DESC`
	SELECT_COMMENTS_BY_UID                 	 = `SELECT c.comment_id, p.post_id, p.title, c.comment_body, c.created_at, c.updated_at 
												FROM comment c 
												JOIN blog_posts p 
												ON c.post_id = p.post_id 
												WHERE uid = ? 
												ORDER BY created_at DESC`
	SELECT_POST_OF_LATEST_UPDATED_CATEGORY 	 = `SELECT post_id, banner, title, body, view_count, comment_count, created_at, updated_at, author, category_id, category_name 
												FROM homepage_latest`
	SELECT_EDITOR_PICK                     	 = `SELECT b.post_id, b.banner, b.title, b.body, b.view_count, b.comment_count, b.created_at, b.updated_at, CONCAT(u.first_name, ' ', u.last_name) AS author 
												FROM blog_posts b 
												JOIN user_details u 
												ON b.author_id = u.uid 
												JOIN editor_pick e 
												ON b.post_id = e.post_id;`
	SELECT_COUNT_OF_CATEGORY_ASSOC_OF_POST   = `SELECT COUNT(*) 
												FROM category_associations 
												WHERE post_id = ?`

	PICK_HEADER_POST = "CALL NEW_EDITOR_PICK(?)"
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

func (p PostRepositoryImpl) IncreaseView(ctx context.Context, id int64) error {
	result, err := p.db.ExecContext(ctx, ADD_VIEW, id)
	if err != nil {
		log.Println("[ERROR] IncreaseView -> error on executing query :", err)
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

func (p PostRepositoryImpl) GetRawPost(ctx context.Context, id int64) (entity.BlogPost, error) {
	var post entity.BlogPost

	rows, err := p.db.QueryContext(ctx, SELECT_RAW_POST, id)
	if err != nil {
		log.Println("[ERROR] GetRawPost -> error on executing query :", err)
		return post, err
	}

	if rows.Next() {
		err = rows.Scan(&post.PostID, &post.AuthorID, &post.Banner, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			log.Println("[ERROR] GetRawPost -> error scanning row :", err)
			return post, err
		}

		return post, nil
	}

	return post, fmt.Errorf("No post with id %d", id)
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
		err = rows.Scan(&post.PostID, &post.AuthorID, &post.Banner, &post.Title, &post.Body, &post.ViewCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt, &author.UserID, &author.Email, &author.FirstName, &author.LastName, &author.Picture, &author.Phone, &author.AboutMe, &author.CreatedAt, &author.UpdatedAt)
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

	rows, err := p.db.QueryContext(ctx, SELECT_LIST_OF_POST, offset, limit)
	if err != nil {
		log.Println("[ERROR] GetBriefsBlogPostData -> error on executing query :", err)
		return postList, err
	}

	for rows.Next() {
		var post entity.BriefBlogPost
		err = rows.Scan(&post.PostID, &post.Banner, &post.Title, &post.Body, &post.ViewCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			log.Println("[ERROR] GetBriefsBlogPostData -> error scanning row :", err)
			return postList, err
		}

		postList = append(postList, &post)
	}

	return postList, nil
}

func (p PostRepositoryImpl) GetMiniBlogPostsDataByUser(ctx context.Context, id int64, limit int64) (entity.PostsTitleWithCategory, error) {
	var postList entity.PostsTitleWithCategory
	var args []interface{}

	query := SELECT_LIST_OF_POST_BY_USER
	args = append(args, id, id)

	if limit > 0 {
		query = query + " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println("[ERROR] GetMiniBlogPostsDataByUser -> error on executing query :", err)
		return postList, err
	}

	for rows.Next() {
		var post entity.PostTitleWithCategory
		err = rows.Scan(&post.PostID, &post.Title, &post.CreatedAt, &post.Category)
		if err != nil {
			log.Println("[ERROR] GetMiniBlogPostsDataByUser -> error scanning row :", err)
			return postList, err
		}

		postList = append(postList, &post)
	}

	return postList, nil
}

func (p PostRepositoryImpl) GetBriefsBlogPostDataOfCategories(ctx context.Context, categories string, offset int64, limit int64) (entity.BriefsBlogPost, error) {
	var postList entity.BriefsBlogPost

	rows, err := p.db.QueryContext(ctx, SELECT_LISF_OF_POST_IN_CATEGORY, categories, offset, limit)
	if err != nil {
		log.Println("[ERROR] GetBriefsBlogPostDataOfCategories -> error on executing query :", err)
		return postList, err
	}

	for rows.Next() {
		var post entity.BriefBlogPost
		err = rows.Scan(&post.PostID, &post.Banner, &post.Title, &post.Body, &post.ViewCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			log.Println("[ERROR] GetBriefsBlogPostDataOfCategories -> error scanning row :", err)
			return postList, err
		}

		postList = append(postList, &post)
	}

	return postList, nil
}

func (p PostRepositoryImpl) GetEditorsPick(ctx context.Context) (entity.BriefBlogPost, error) {
	var post entity.BriefBlogPost

	rows, err := p.db.QueryContext(ctx, SELECT_EDITOR_PICK)
	if err != nil {
		log.Println("[ERROR] GetEditorPick -> error on executing query :", err)
		return post, err
	}

	if rows.Next() {
		err = rows.Scan(&post.PostID, &post.Banner, &post.Title, &post.Body, &post.ViewCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			log.Println("[ERROR] GetEditorPick -> error scanning row :", err)
			return post, err
		}

		return post, nil
	}

	return post, errors.New("can't get editor pick post")
}

func (p PostRepositoryImpl) GetTopCategoryPost(ctx context.Context) (entity.BriefsBlogPost, entity.Categories, error) {
	var postList entity.BriefsBlogPost
	var categoryList entity.Categories

	rows, err := p.db.QueryContext(ctx, SELECT_POST_OF_LATEST_UPDATED_CATEGORY)
	if err != nil {
		log.Println("[ERROR] GetTopCategoryPost -> error on executing query :", err)
		return postList, categoryList, err
	}

	for rows.Next() {
		var post entity.BriefBlogPost
		var category entity.Category
		err = rows.Scan(&post.PostID, &post.Banner, &post.Title, &post.Body, &post.ViewCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt, &post.Author, &category.CategoryID, &category.CategoryName)
		if err != nil {
			log.Println("[ERROR] GetTopCategoryPost -> error scanning row :", err)
			return postList, categoryList, err
		}

		postList = append(postList, &post)
		categoryList = append(categoryList, &category)
	}

	return postList, categoryList, nil
}

func (p PostRepositoryImpl) GetBriefsBlogPostFromSearch(ctx context.Context, q string, offset int64, limit int64, dateStart *time.Time, dateEnd *time.Time, category string) (entity.BriefsBlogPost, error) {
	var postList entity.BriefsBlogPost
	var query string
	var args []interface{}

	args = append(args, fmt.Sprintf("%%%s%%", q))
	query = SELECT_FULL_TEXT_POST

	if category != "" {
		query = query + " JOIN category_associations a ON  a.post_id = b.post_id JOIN categories c ON c.category_id = a.category_id"
	}

	// Add where clause
	query = query + " WHERE b.title LIKE ?"

	// Add argument
	if category != "" {
		query = query + " AND c.category_name = ?"
		args = append(args, category)
	}

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
		err = rows.Scan(&post.PostID, &post.Banner, &post.Title, &post.Body, &post.ViewCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			log.Println("[ERROR] GetBriefsBlogPostFromSearch -> error scanning row :", err)
			return postList, err
		}

		postList = append(postList, &post)
	}

	return postList, nil
}

func (p PostRepositoryImpl) CountSearchResult(ctx context.Context, q string, dateStart *time.Time, dateEnd *time.Time, category string) (int64, error) {
	var count int64
	var query string
	var args []interface{}

	args = append(args, fmt.Sprintf("%%%s%%", q))
	query = COUNT_SEARCH_RESULT
	if category != "" {
		query = query + " JOIN category_associations a ON  a.post_id = b.post_id JOIN categories c ON c.category_id = a.category_id"
	}

	// Add where clause
	query = query + " WHERE b.title LIKE ?"

	// Add argument
	if category != "" {
		query = query + " AND c.category_name = ?"
		args = append(args, category)
	}

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

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println("[ERROR] CountSearchResult -> error on executing query :", err)
		return 0, err
	}

	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println("[ERROR] CountSearchResult -> error scanning row :", err)
			return 0, err
		}

		return count, nil
	}

	return 0, errors.New("can't get count of post ")
}

func (p PostRepositoryImpl) GetUserPostStatisticOfEachCategory(ctx context.Context, id int64) (entity.ListOfUserPostStatisticByCategory, error) {
	var stats entity.ListOfUserPostStatisticByCategory

	rows, err := p.db.QueryContext(ctx, SELECT_EACH_CATEGORY_USER_POST_STATISTIC, id)
	if err != nil {
		log.Println("[ERROR] GetUserPostStatisticOfEachCategory -> error on executing query :", err)
		return stats, err
	}

	for rows.Next() {
		var stat entity.UserPostStatisticByCategory

		err = rows.Scan(&stat.Category, &stat.TotalPost, &stat.TotalView)
		if err != nil {
			log.Println("[ERROR] GetUserPostStatisticOfEachCategory -> error scanning row :", err)
			return stats, err
		}
		stats = append(stats, &stat)
	}

	return stats, nil
}

func (p PostRepositoryImpl) CountListOfPost(ctx context.Context) (int64, error) {
	var count int64

	rows, err := p.db.QueryContext(ctx, COUNT_LIST_OF_POST)
	if err != nil {
		log.Println("[ERROR] CountListOfPost -> error on executing query :", err)
		return 0, err
	}

	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println("[ERROR] CountListOfPost -> error scanning row :", err)
			return 0, err
		}

		return count, nil
	}

	return 0, errors.New("can't get count of post ")
}

func (p PostRepositoryImpl) CountUserTotalPost(ctx context.Context, id int64) (int64, error) {
	var count int64

	rows, err := p.db.QueryContext(ctx, COUNT_TOTAL_USER_POST, id)
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

func (p PostRepositoryImpl) CountListOfPostInCategories(ctx context.Context, categories string) (int64, error) {
	var count int64

	rows, err := p.db.QueryContext(ctx, COUNT_LIST_OF_POST_IN_CATEGOIES, categories)
	if err != nil {
		log.Println("[ERROR] CountListOfPostInCategories -> error on executing query :", err)
		return 0, err
	}

	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println("[ERROR] CountListOfPostInCategories -> error scanning row :", err)
			return 0, err
		}

		return count, nil
	}

	return 0, errors.New("can't get count of post ")
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
		categories = append(categories, &entity.Category{
			CategoryID:   -1,
			CategoryName: "No Categories",
		})
		return categories, nil
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

func (p PostRepositoryImpl) AddPostCategoryAssoc(ctx context.Context, postId int64, categoryId int64) error {
	result, err := p.db.ExecContext(ctx, INSERT_CATEGORY_ASSOC, postId, categoryId)
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

func (p PostRepositoryImpl) UpsertPostCategoryAssoc(ctx context.Context, postId int64, categoryId int64) error {
	rows, err := p.db.QueryContext(ctx, SELECT_COUNT_OF_CATEGORY_ASSOC_OF_POST, postId)
	if err != nil {
		log.Println("[ERROR] UpsertPostCategoryAssoc -> error on executing query :", err)
		return err
	}

	var count int64
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println("[ERROR] UpsertPostCategoryAssoc -> error scanning row :", err)
			return err
		}
	} else {
		return errors.New("can't count how many record on category assoc")
	}

	if count == 0 {
		_, err = p.db.ExecContext(ctx, INSERT_CATEGORY_ASSOC, postId, categoryId)
		if err != nil {
			log.Println("[ERROR] UpsertPostCategoryAssoc -> error inserting row")
			return err
		}

		return nil
	}

	_, err = p.db.ExecContext(ctx, UPDATE_CATEGORY_ASSOC, categoryId, postId)
	if err != nil {
		log.Println("[ERROR] UpsertPostCategoryAssoc -> error inserting row :", err)
		return err
	}
	return nil
}

func (p PostRepositoryImpl) DeletePostCategoryAssoc(ctx context.Context, postId int64) error {
	result, err := p.db.ExecContext(ctx, DELETE_CATEGORY_ASSOC, postId)
	if err != nil {
		log.Println("[ERROR] DeletePostCategoryAssoc -> error inserting row :", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] DeletePostCategoryAssoc -> error on getting rows affected :", err)
		return err
	}
	if rowsAffected == 0 {
		log.Println("[ERROR] DeletePostCategoryAssoc -> error on updating row :", err)
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

func (p PostRepositoryImpl) DeleteComment(ctx context.Context, id int64) error {
	result, err := p.db.ExecContext(ctx, DELETE_COMMENT, id)
	if err != nil {
		log.Println("[ERROR] DeleteComment -> error deleting row :", err)
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

func (p PostRepositoryImpl) GetPostComments(ctx context.Context, id int64) (entity.Comments, entity.MiniUsersDetail, error) {
	var comments entity.Comments
	var users entity.MiniUsersDetail

	rows, err := p.db.QueryContext(ctx, SELECT_COMMENTS, id)
	if err != nil {
		log.Println("[ERROR] GetPostComments -> error on executing query :", err)
		return comments, users, err
	}

	for rows.Next() {
		var comment entity.Comment
		var user entity.MiniUserDetail
		err = rows.Scan(&comment.CommentID, &comment.UserID, &comment.CommentBody, &comment.CreatedAt, &comment.UpdatedAt, &user.UserID, &user.FirstName, &user.LastName, &user.Picture)
		if err != nil {
			log.Println("[ERROR] GetPostComments -> error scanning row :", err)
			return comments, users, err
		}

		comments = append(comments, &comment)
		users = append(users, &user)
	}

	return comments, users, nil
}

func (p PostRepositoryImpl) GetCommentAuthorId(ctx context.Context, id int64) (int64, error) {
	rows, err := p.db.QueryContext(ctx, SELECT_COMMENT_AUTHOR, id)
	if err != nil {
		log.Println("[ERROR] GetCommentAuthorId -> error on executing query :", err)
		return -1, err
	}

	var uid int64
	if rows.Next() {
		err = rows.Scan(&uid)
		if err != nil {
			log.Println("[ERROR] GetCommentAuthorId -> error scanning row :", err)
			return -1, err
		}

		return uid, nil
	}

	return -1, fmt.Errorf("No comment with id %d", id)
}

func (p PostRepositoryImpl) GetUserComments(ctx context.Context, id int64, limit int64) (entity.Comments, entity.BriefsBlogPost, error) {
	var comments entity.Comments
	var posts entity.BriefsBlogPost
	var args []interface{}

	query := SELECT_COMMENTS_BY_UID
	args = append(args, id)

	if limit > 0 {
		query = query + " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println("[ERROR] GetUserComments -> error on executing query :", err)
		return comments, posts, err
	}

	for rows.Next() {
		var comment entity.Comment
		var post entity.BriefBlogPost
		err = rows.Scan(&comment.CommentID, &post.PostID, &post.Title, &comment.CommentBody, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			log.Println("[ERROR] GetPostComments -> error scanning row :", err)
			return comments, posts, err
		}

		comments = append(comments, &comment)
		posts = append(posts, &post)
	}

	return comments, posts, nil
}

func (p PostRepositoryImpl) CountUserTotalComment(ctx context.Context, id int64) (int64, error) {
	var count int64

	rows, err := p.db.QueryContext(ctx, COUNT_TOTAL_USER_COMMENT, id)
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

func (p PostRepositoryImpl) PicHeaderPost(ctx context.Context, postId int64) error {
	result, err := p.db.ExecContext(ctx, PICK_HEADER_POST, postId)
	if err != nil {
		log.Println("[ERROR] PicHeaderPost -> error inserting row :", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] PicHeaderPost -> error on getting rows affected :", err)
		return err
	}
	if rowsAffected != 1 {
		log.Println("[ERROR] PicHeaderPost -> error on updating row :", err)
		return err
	}

	return nil
}