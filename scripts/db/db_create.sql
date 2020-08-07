CREATE TABLE setting (
  setting_key VARCHAR(100) NOT NULL,
  setting_value VARCHAR(100),
  PRIMARY KEY (setting_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO setting (setting_key, setting_value)
  VALUES ('db_version', '1');

CREATE TABLE role (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(20) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO role (id, name)
  VALUES (1, 'admin'), (2, 'evaluator'), (3, 'user');

CREATE TABLE user (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  username VARCHAR(100) NOT NULL,
  password VARCHAR(100) NOT NULL,
  must_change_password TINYINT(1) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE user_role (
  user_id INT NOT NULL,
  role_id INT NOT NULL,
  PRIMARY KEY (user_id, role_id),
  KEY fk_userrole_user (user_id),
  KEY fk_userrole_role (role_id),
  CONSTRAINT fk_userrole_user FOREIGN KEY (user_id)
    REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_userrole_role FOREIGN KEY (role_id)
    REFERENCES role (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE user_setting (
  user_id INT NOT NULL,
  setting_key VARCHAR(100) NOT NULL,
  setting_value VARCHAR(1000) NOT NULL,
  PRIMARY KEY (user_id, setting_key),
  KEY fk_usersetting_user (user_id),
  CONSTRAINT fk_usersetting_user FOREIGN KEY (user_id)
    REFERENCES user (id) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE user_contract (
  user_id INT NOT NULL,
  daily_working_duration INT NOT NULL,
  annual_vacation_days FLOAT NOT NULL,
  init_overtime_duration INT NOT NULL,
  init_vacation_days FLOAT NOT NULL,
  first_work_day TIMESTAMP NOT NULL,
  PRIMARY KEY (user_id),
  KEY fk_usercontract_user (user_id),
  CONSTRAINT fk_usercontract_user FOREIGN KEY (user_id)
    REFERENCES user (id) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE session (
  id VARCHAR(32) NOT NULL,
  user_id INT,
  expire_at TIMESTAMP NOT NULL,
  previous_url VARCHAR(100),
  PRIMARY KEY (id),
  KEY fk_session_user (user_id),
  CONSTRAINT fk_session_user FOREIGN KEY (user_id)
    REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE INDEX idx_session_expireat ON session(expire_at);

CREATE TABLE entry_type (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO entry_type (id, name)
  VALUES (1, 'work'), (2, 'travel'), (3, 'vacation'), (4, 'holiday'), (5, 'illness');

CREATE TABLE entry_activity (
  id INT NOT NULL AUTO_INCREMENT,
  description VARCHAR(50) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE entry (
  id INT NOT NULL AUTO_INCREMENT,
  user_id INT NOT NULL,
  type_id INT NOT NULL,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL,
  break_duration INT NOT NULL,
  activity_id INT,
  description VARCHAR(200),
  PRIMARY KEY (id),
  KEY fk_entry_user (user_id),
  KEY fk_entry_entrytype (type_id),
  KEY fk_entry_entryactivity (activity_id),
  CONSTRAINT fk_entry_user FOREIGN KEY (user_id)
    REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_entry_entrytype FOREIGN KEY (type_id)
    REFERENCES entry_type (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT fk_entry_entryactivity FOREIGN KEY (activity_id)
    REFERENCES entry_activity (id) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE INDEX idx_entry_starttime ON entry(start_time);
CREATE INDEX idx_entry_endtime ON entry(end_time);