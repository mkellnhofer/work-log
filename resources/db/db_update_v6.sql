CREATE TABLE project (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(30) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY unique_project_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE entry ADD COLUMN project_id INT DEFAULT NULL AFTER activity_id;

ALTER TABLE entry ADD KEY fk_entry_project (project_id);

ALTER TABLE entry ADD CONSTRAINT fk_entry_project FOREIGN KEY (project_id)
  REFERENCES project (id) ON DELETE SET NULL ON UPDATE CASCADE;