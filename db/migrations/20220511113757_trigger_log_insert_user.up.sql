CREATE TRIGGER IF NOT EXISTS log_insert_users
AFTER INSERT ON users
FOR EACH ROW 
BEGIN
	IF (SELECT COUNT(logger.log_id) > 200 FROM logger)
	THEN 
		DELETE FROM logger
		WHERE log_id = (SELECT MAX(log_id)FROM logger) - 200;    
	END IF;
	INSERT INTO logger(operation, statement, table_name)
	VALUE ("Insert", CONCAT("User with ID ", NEW.uid, " is created"), "users");
END