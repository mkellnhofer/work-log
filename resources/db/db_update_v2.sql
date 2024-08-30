CREATE TABLE contract (
  user_id INT NOT NULL,
  init_overtime_hours FLOAT NOT NULL,
  init_vacation_days FLOAT NOT NULL,
  first_day DATE NOT NULL DEFAULT '0000-00-00',
  PRIMARY KEY (user_id),
  KEY fk_contract_user (user_id),
  CONSTRAINT fk_contract_user FOREIGN KEY (user_id)
    REFERENCES user (id) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE contract_working_hours(
  user_id INT NOT NULL,
  first_day DATE NOT NULL DEFAULT '0000-00-00',
  daily_hours FLOAT NOT NULL,
  PRIMARY KEY (user_id, first_day),
  KEY fk_contractworkinghours_contract (user_id),
  CONSTRAINT fk_contractworkinghours_contract FOREIGN KEY (user_id)
    REFERENCES contract (user_id) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE contract_vacation_days (
  user_id INT NOT NULL,
  first_day DATE NOT NULL DEFAULT '0000-00-00',
  monthly_days FLOAT NOT NULL,
  PRIMARY KEY (user_id, first_day),
  KEY fk_contractvacationdays_contract (user_id),
  CONSTRAINT fk_contractvacationdays_contract FOREIGN KEY (user_id)
    REFERENCES contract (user_id) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO contract (user_id, init_overtime_hours, init_vacation_days, first_day)
  SELECT user_id, init_overtime_duration / 60, init_vacation_days, DATE(first_work_day) FROM user_contract;

INSERT INTO contract_working_hours (user_id, first_day, daily_hours)
  SELECT user_id, DATE(first_work_day), daily_working_duration / 60 FROM user_contract;

INSERT INTO contract_vacation_days (user_id, first_day, monthly_days)
  SELECT user_id, DATE(first_work_day), annual_vacation_days / 12 FROM user_contract;

DROP TABLE user_contract;