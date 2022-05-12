CREATE TRIGGER IF NOT EXISTS log_delete_post
AFTER DELETE ON blog_posts
FOR EACH ROW 
BEGIN
	DECLARE exc_query VARCHAR(1024);
	SET exc_query = (SELECT info FROM INFORMATION_SCHEMA.PROCESSLIST
					 WHERE id = CONNECTION_ID());
	
	IF (SELECT COUNT(logger.log_id) > 200 FROM logger)
	THEN 
		DELETE FROM logger
		WHERE log_id = (SELECT MAX(log_id)FROM logger) - 200;    
	END IF;
	INSERT INTO logger(operation, statement, tabel)
	VALUE ("Delete", exc_query, "blog_post");
END