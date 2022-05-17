CREATE TRIGGER update_latest_post_of_category
AFTER INSERT ON category_associations
FOR EACH ROW 
BEGIN
	UPDATE categories
    SET latest_post = CURRENT_TIMESTAMP
    WHERE category_id = NEW.category_id;
END