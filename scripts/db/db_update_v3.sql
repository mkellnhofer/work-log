UPDATE `entry` SET end_time = DATE_SUB(end_time, INTERVAL break_duration MINUTE);

ALTER TABLE `entry` DROP COLUMN break_duration;