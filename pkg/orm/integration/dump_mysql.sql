/*Generated by orm 2022-01-31 00:56:48, from sqlite3 to mysql*/

SET sql_mode='NO_BACKSLASH_ESCAPES';
CREATE TABLE IF NOT EXISTS `test_dump_struct` (`id` INTEGER PRIMARY KEY AUTO_INCREMENT NOT NULL, `name` TEXT NULL, `is_man` INTEGER NULL, `created` DATETIME NULL);
INSERT INTO `test_dump_struct` (`id`, `name`, `is_man`, `created`) VALUES (1,'1',1,'2022-01-30T19:26:48Z');
INSERT INTO `test_dump_struct` (`id`, `name`, `is_man`, `created`) VALUES (2,CONCAT('2', CHAR(10)),0,'2022-01-30T19:26:48Z');
INSERT INTO `test_dump_struct` (`id`, `name`, `is_man`, `created`) VALUES (3,'3;',0,'2022-01-30T19:26:48Z');
INSERT INTO `test_dump_struct` (`id`, `name`, `is_man`, `created`) VALUES (4,CONCAT('4', CHAR(10), ';', CHAR(10), ''''''),0,'2022-01-30T19:26:48Z');
INSERT INTO `test_dump_struct` (`id`, `name`, `is_man`, `created`) VALUES (5,CONCAT('5''', CHAR(10)),0,'2022-01-30T19:26:48Z');