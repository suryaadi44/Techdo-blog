CREATE VIEW homepage_latest AS
SELECT * FROM (
	SELECT b.post_id, b.banner, b.title, b.body, b.view_count, b.comment_count, b.created_at, b.updated_at, CONCAT(u.first_name, ' ', u.last_name) AS author, 
		c.category_id, c.category_name, DENSE_RANK() OVER (PARTITION BY c.category_id ORDER BY b.created_at DESC) r
	FROM category_associations a
		JOIN blog_posts b ON b.post_id = a.post_id
        JOIN user_details u ON b.author_id = u.uid 
		JOIN (
			SELECT category_id, category_name FROM categories 
			ORDER BY latest_post DESC LIMIT 3
		) c ON a.category_id IN (c.category_id)
) AS sub WHERE sub.r <= 3 ORDER BY sub.created_at DESC;