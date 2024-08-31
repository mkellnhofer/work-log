-- ### Test data for development and manual testing ###
-- ( Uses relative dates (based on CURDATE()) so data stays current. )

-- Date anchors
SET @contract_start = DATE_SUB(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL 2 MONTH);
SET @contract_change = DATE_SUB(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL 1 MONTH);
SET @w1_mon = DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) + 7 DAY);
SET @w1_tue = @w1_mon + INTERVAL 1 DAY;
SET @w1_wed = @w1_mon + INTERVAL 2 DAY;
SET @w1_thu = @w1_mon + INTERVAL 3 DAY;

-- Clear existing data
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE entry_label;
TRUNCATE TABLE label;
TRUNCATE TABLE entry;
TRUNCATE TABLE project;
TRUNCATE TABLE entry_activity;
TRUNCATE TABLE contract_vacation_days;
TRUNCATE TABLE contract_working_hours;
TRUNCATE TABLE contract;
TRUNCATE TABLE session;
TRUNCATE TABLE user_setting;
TRUNCATE TABLE user_role;
TRUNCATE TABLE user;
SET FOREIGN_KEY_CHECKS = 1;

-- User ("john.doe"/"password") (password must be changed on first login)
INSERT INTO user (id, name, username, password, must_change_password)
  VALUES
    (1, 'John Doe', 'john.doe', '$2a$10$/nwF.qQuDWCa.7e2uCiZHu.X.4Bav3Gj9ebCHHmcJcp34lS61TPli', 1);
INSERT INTO user_role (user_id, role_id)
  VALUES
    (1, 1), (1, 2), (1, 3);

-- Contract (starts 2 months ago, changes 1 month ago)
INSERT INTO contract (user_id, init_overtime_hours, init_vacation_days, first_day)
  VALUES
    (1, 300.0, 3.5, @contract_start);
INSERT INTO contract_vacation_days (user_id, first_day, monthly_days)
  VALUES
    (1, @contract_start, 2.0),
    (1, @contract_change, 2.5);
INSERT INTO contract_working_hours (user_id, first_day, daily_hours)
  VALUES
    (1, @contract_start, 8.0),
    (1, @contract_change, 7.5);

-- Entry activities
INSERT INTO entry_activity (id, description)
  VALUES
    (1, 'General'),
    (2, 'Meeting'),
    (3, 'Organization'),
    (4, 'Research'),
    (5, 'Training'),
    (6, 'Conception'),
    (7, 'Development'),
    (8, 'Testing'),
    (9, 'Documentation'),
    (10, 'Trouble Shooting'),
    (11, 'Infrastructure'),
    (12, 'Help'),
    (13, 'Support');

-- Projects
INSERT INTO project (id, name)
  VALUES
    (1, 'Backend'),
    (2, 'Web App');

-- Entries (anchored to last week's Mon-Thu)
INSERT INTO entry (id, user_id, type_id, start_time, end_time, activity_id, project_id, description)
  VALUES
    (1,  1, 1, TIMESTAMP(@w1_mon, '08:00:00'), TIMESTAMP(@w1_mon, '12:00:00'), 7, 2, 'Overview UI'),
    (2,  1, 1, TIMESTAMP(@w1_mon, '13:00:00'), TIMESTAMP(@w1_mon, '14:00:00'), 4, 1, 'New build system'),
    (3,  1, 1, TIMESTAMP(@w1_mon, '14:00:00'), TIMESTAMP(@w1_mon, '17:30:00'), 7, 2, 'Overview UI'),
    (4,  1, 1, TIMESTAMP(@w1_tue, '09:25:00'), TIMESTAMP(@w1_tue, '09:55:00'), 1, NULL, 'Watering office flowers'),
    (5,  1, 1, TIMESTAMP(@w1_tue, '09:55:00'), TIMESTAMP(@w1_tue, '12:00:00'), 6, 1, 'Migration micro services'),
    (6,  1, 1, TIMESTAMP(@w1_tue, '13:00:00'), TIMESTAMP(@w1_tue, '14:30:00'), 2, 1, 'New build system'),
    (7,  1, 1, TIMESTAMP(@w1_tue, '14:15:00'), TIMESTAMP(@w1_tue, '18:30:00'), 7, 2, 'Overview UI'),
    (8,  1, 3, TIMESTAMP(@w1_wed, '09:00:00'), TIMESTAMP(@w1_wed, '17:00:00'), NULL, NULL, NULL),
    (9,  1, 1, TIMESTAMP(@w1_thu, '09:00:00'), TIMESTAMP(@w1_thu, '09:35:00'), 2, NULL, 'Daily standup'),
    (10, 1, 1, TIMESTAMP(@w1_thu, '09:35:00'), TIMESTAMP(@w1_thu, '09:45:00'), 7, 2, 'Bugfix'),
    (11, 1, 1, TIMESTAMP(@w1_thu, '09:45:00'), TIMESTAMP(@w1_thu, '12:15:00'), 7, 2, 'Overview UI'),
    (12, 1, 1, TIMESTAMP(@w1_thu, '13:00:00'), TIMESTAMP(@w1_thu, '14:00:00'), 3, 2, 'Release');

-- Labels
INSERT INTO label (id, name)
  VALUES
    (1, '#WA-247'),
    (2, '#WA-255'),
    (3, 'Bug'),
    (4, 'v5.7.0');

INSERT INTO entry_label (entry_id, label_id)
  VALUES
    (1, 2),
    (2, 2),
    (7, 2),
    (10, 1),
    (10, 3),
    (11, 2),
    (12, 4);
