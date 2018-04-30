CREATE TABLE IF NOT EXISTS `users` (
  `uuid` int(11) NOT NULL AUTO_INCREMENT,
  `displayName` varchar(256) COLLATE utf8_bin NOT NULL,
  `name` varchar(128) COLLATE utf8_bin NOT NULL,
  `password` char(128) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
COMMIT;


CREATE TABLE IF NOT EXISTS `projects` (
  `uuid` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(256) COLLATE utf8_bin NOT NULL,
  `description` varchar(1024) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
COMMIT;


CREATE TABLE IF NOT EXISTS `users_in_projects` (
  `user_id` int(11) NOT NULL,
  `project_id` int(11) NOT NULL,
  `enable` tinyint(1) NOT NULL,
  `expiration` datetime DEFAULT NULL,
  KEY `project_id` (`project_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

ALTER TABLE `users_in_projects`
  ADD CONSTRAINT `users_in_projects_ibfk_2` FOREIGN KEY (`project_id`) REFERENCES `projects` (`uuid`),
  ADD CONSTRAINT `users_in_projects_ibfk_3` FOREIGN KEY (`user_id`) REFERENCES `users` (`uuid`);
COMMIT;