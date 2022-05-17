CREATE TRIGGER increase_comment
AFTER INSERT ON comment
FOR EACH ROW 
BEGIN
	UPDATE blog_posts
    SET comment_count = comment_count + 1
    WHERE post_id = NEW.post_id;
END