CREATE TRIGGER decrease_comment
AFTER DELETE ON comment
FOR EACH ROW 
BEGIN
	UPDATE blog_posts
    SET comment_count = comment_count - 1
    WHERE post_id = OLD.post_id;
END