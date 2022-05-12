CREATE TRIGGER IF NOT EXISTS log_update_post
AFTER UPDATE ON blog_posts
FOR EACH ROW 
BEGIN
	IF (SELECT COUNT(logger.log_id) > 200 FROM logger)
	THEN 
		DELETE FROM logger
		WHERE log_id = (SELECT MAX(log_id)FROM logger) - 200;
	END IF;
	INSERT INTO logger(operation, statement, table_name)
	VALUE ("Update", CONCAT("Post with ID ", OLD.post_id, " is updated"), "blog_post");
END