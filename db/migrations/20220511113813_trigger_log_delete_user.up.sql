CREATE TRIGGER log_delete_users
AFTER DELETE ON users
FOR EACH ROW 
BEGIN
	IF (SELECT COUNT(logger.log_id) > 200 FROM logger)
	THEN 
		DELETE FROM logger
		WHERE log_id = (SELECT MAX(log_id)FROM logger) - 200;    
	END IF;
	INSERT INTO logger(operation, statement, table_name)
	VALUE ("Delete", CONCAT("User with ID ", OLD.uid, " is deleted"), "users");
END