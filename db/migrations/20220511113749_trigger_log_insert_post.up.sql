CREATE TRIGGER log_insert_post
AFTER INSERT ON blog_posts
FOR EACH ROW 
BEGIN
	IF (SELECT COUNT(logger.log_id) > 200 FROM logger)
	THEN 
		DELETE FROM logger
		WHERE log_id = (SELECT MAX(log_id)FROM logger) - 200;    
	END IF;
	INSERT INTO logger(operation, statement, table_name)
	VALUE ("Insert", CONCAT("Post with ID ", NEW.post_id, " is created"), "blog_post");
END