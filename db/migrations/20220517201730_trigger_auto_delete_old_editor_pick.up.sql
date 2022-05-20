CREATE TRIGGER delete_old_editor_pick
AFTER INSERT ON editor_pick
FOR EACH ROW 
BEGIN
    IF (SELECT COUNT(*) > 1 FROM editor_pick)
	THEN 
		DELETE FROM editor_pick
        WHERE id != NEW.id;
	END IF;
END