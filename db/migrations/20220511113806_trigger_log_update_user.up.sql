CREATE TRIGGER log_update_users
BEFORE UPDATE ON users
FOR EACH ROW 
BEGIN
	IF (SELECT COUNT(logger.log_id) > 200 FROM logger)
	THEN 
		DELETE FROM logger
		WHERE log_id = (SELECT MAX(log_id)FROM logger) - 200;    
	END IF;
	INSERT INTO logger(operation, statement, table_name)
	VALUE ("Update", CONCAT("User with ID ", OLD.uid, " is updated"), "users");
END