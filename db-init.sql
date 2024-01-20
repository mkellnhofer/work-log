CREATE DATABASE work_log;

CREATE USER 'work_log'@'%' IDENTIFIED BY 'work_log';
GRANT ALL ON work_log.* TO 'work_log'@'%' WITH GRANT OPTION;