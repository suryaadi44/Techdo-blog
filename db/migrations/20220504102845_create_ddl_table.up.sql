CREATE TABLE IF NOT EXISTS users (
	uid INT NOT NULL AUTO_INCREMENT,
	username VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    PRIMARY KEY(uid)
)ENGINE=InnoDB;


CREATE TABLE IF NOT EXISTS user_details(
	uid INT NOT NULL,
    email VARCHAR(25) NOT NULL,
    first_name VARCHAR(20) NOT NULL,
    last_name VARCHAR(20) NOT NULL,
    picture VARCHAR(50) NOT NULL,
    phone VARCHAR(15) DEFAULT "",
    about_me TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(uid)
)ENGINE=InnoDB;


CREATE TABLE IF NOT EXISTS sessions(
	token VARCHAR(40) NOT NULL,
	uid INT NOT NULL, 
	expireAt DATETIME,
	PRIMARY KEY (token),
	FOREIGN KEY (uid) REFERENCES users(uid)
)ENGINE=InnoDB;


CREATE TABLE IF NOT EXISTS blog_posts(
	post_id INT NOT NULL AUTO_INCREMENT,
    author_id INT DEFAULT NULL,
    banner VARCHAR(100) DEFAULT NULL,
    title VARCHAR(100) DEFAULT NULL,
    body TEXT DEFAULT NULL,
<<<<<<< HEAD
	view_count INT DEFAULT 0,
	comment_count INT DEFAULT 0,
=======
    view_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
>>>>>>> 8ad429016ef56b60dc2df0bba0c97d5c64b8428a
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id),
    FOREIGN KEY (author_id) REFERENCES users(uid)
)ENGINE=InnoDB;


CREATE TABLE IF NOT EXISTS categories(
	category_id INT NOT NULL AUTO_INCREMENT,
    category_name VARCHAR(25) NOT NULL,
    PRIMARY KEY (category_id)
)ENGINE=InnoDB;


CREATE TABLE IF NOT EXISTS category_associations(
	post_id INT NOT NULL,
    category_id INT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES blog_posts(post_id),
    FOREIGN KEY (category_id) REFERENCES categories(category_id)
)ENGINE=InnoDB;


CREATE TABLE IF NOT EXISTS comment(
	comment_id INT NOT NULL AUTO_INCREMENT,
    post_id INT NOT NULL,
    uid INT NOT NULL,
    comment_body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (comment_id),
    FOREIGN KEY (post_id) REFERENCES blog_posts(post_id),
    FOREIGN KEY (uid) REFERENCES users(uid)
)ENGINE=InnoDB;	

CREATE TABLE IF NOT EXISTS logger(
	log_id INT AUTO_INCREMENT, 
	operation VARCHAR(6),
	statement TEXT,
	table_name VARCHAR(9),
	time_stamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (log_id)
)ENGINE=InnoDB;
