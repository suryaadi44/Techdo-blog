CREATE TABLE IF NOT EXISTS users (
	uid INT NOT NULL AUTO_INCREMENT,
	username VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL
    PRIMARY KEY(uid)
)ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS user_details(
	uid INT NOT NULL,
    email VARCHAR(25) NOT NULL,
    first_name VARCHAR(20) NOT NULL,
    last_name VARCHAR(20),
    picture VARCHAR(50),
    phone VARCHAR(15),
    about_me TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(uid)
)ENGINE=InnoDB;
CREATE TABLE sessions(
	token INT NOT NULL,
	uid INT NOT NULL, 
	expireAt DATETIME,
	PRIMARY KEY (token),
	FOREIGN KEY (uid) REFERENCES users(uid)
)
CREATE TABLE blog_posts(
	post_id INT NOT NULL AUTO_INCREMENT,
    author_id INT UNIQUE,
    banner VARCHAR(50) NOT NULL,
    title VARCHAR(50) NOT NULL,
    body TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id),
    FOREIGN KEY (author_id) REFERENCES users(uid)
)ENGINE=InnoDB;
CREATE TABLE categories(
	category_id INT NOT NULL AUTO_INCREMENT,
    category_name VARCHAR(25) NOT NULL,
    PRIMARY KEY (category_id)
)ENGINE=InnoDB;
CREATE TABLE category_associations(
	post_id INT UNIQUE NOT NULL,
    category_id INT UNIQUE NOT NULL,
    FOREIGN KEY (post_id) REFERENCES blog_posts(post_id),
    FOREIGN KEY (category_id) REFERENCES categories(category_id)
)ENGINE=InnoDB;
CREATE TABLE comment(
	comment_id INT NOT NULL AUTO_INCREMENT,
    post_id INT UNIQUE NOT NULL,
    uid INT UNIQUE NOT NULL,
    comment_body TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (comment_id),
    FOREIGN KEY (post_id) REFERENCES blog_posts(post_id),
    FOREIGN KEY (uid) REFERENCES users(uid)
)ENGINE=InnoDB;	