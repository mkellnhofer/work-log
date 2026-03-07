CREATE TABLE token (
  id INT NOT NULL AUTO_INCREMENT,
  user_id INT NOT NULL,
  name VARCHAR(30) NOT NULL,
  hashed_token VARCHAR(64) NOT NULL,
  truncated_token VARCHAR(32) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY unique_token (hashed_token),
  KEY fk_token_user (user_id),
  CONSTRAINT fk_token_user FOREIGN KEY (user_id)
    REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
