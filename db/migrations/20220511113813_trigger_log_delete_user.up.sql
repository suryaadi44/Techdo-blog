CREATE TRIGGER log_delete_users
BEFORE DELETE ON users
FOR EACH ROW 
BEGIN
	DELETE FROM blog_posts 
	WHERE blog_posts.author_id = OLD.uid;
	IF (SELECT COUNT(logger.log_id) > 200 FROM logger)
	THEN 
		DELETE FROM logger
		WHERE log_id = (SELECT MAX(log_id)FROM logger) - 200;    
	END IF;
	INSERT INTO logger(operation, statement, table_name)
	VALUE ("Delete", CONCAT("User with ID ", OLD.uid, " is deleted"), "users");
END