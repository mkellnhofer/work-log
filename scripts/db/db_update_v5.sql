CREATE TABLE label (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(20) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY unique_label_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE entry_label (
  entry_id INT NOT NULL,
  label_id INT NOT NULL,
  PRIMARY KEY (entry_id, label_id),
  KEY fk_entrylabel_entry (entry_id),
  KEY fk_entrylabel_label (label_id),
  CONSTRAINT fk_entrylabel_entry FOREIGN KEY (entry_id)
    REFERENCES entry (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_entrylabel_label FOREIGN KEY (label_id)
    REFERENCES label (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;