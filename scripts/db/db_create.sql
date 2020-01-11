CREATE TABLE setting (
  setting_key VARCHAR(100) NOT NULL,
  setting_value VARCHAR(100),
  PRIMARY KEY (setting_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO setting (setting_key, setting_value)
  VALUES ('db_version', '1');

CREATE TABLE user (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  username VARCHAR(100) NOT NULL,
  password VARCHAR(100) NOT NULL,
  PRIMARY KEY (id)
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
  description VARCHAR(50) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO entry_type (id, description)
  VALUES (1, 'Arbeit'), (2, 'Reise'), (3, 'Urlaub'), (4, 'Feiertag'), (5, 'Krankheit');

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