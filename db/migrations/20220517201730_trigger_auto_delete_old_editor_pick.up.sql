CREATE TRIGGER delete_old_editor_pick
AFTER INSERT ON editor_pick
FOR EACH ROW 
BEGIN
	DELETE FROM editor_pick
    WHERE id != NEW.id;
END