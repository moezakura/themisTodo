CREATE DATABASE IF NOT EXISTS `themis_todo` DEFAULT CHARACTER SET utf8 COLLATE utf8_bin;
USE `themis_todo`;

CREATE TABLE IF NOT EXISTS `authority_role` (
  `user_id` int(11) NOT NULL,
  `type` int(11) NOT NULL,
  `enable` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE IF NOT EXISTS `projects` (
  `uuid` int(11) NOT NULL,
  `name` varchar(256) COLLATE utf8_bin NOT NULL,
  `description` varchar(1024) COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE IF NOT EXISTS `todo_list` (
  `id` int(11) NOT NULL,
  `project` int(11) NOT NULL,
  `name` varchar(1024) COLLATE utf8_bin NOT NULL,
  `creator` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  `deadline` date NOT NULL,
  `description` text COLLATE utf8_bin NOT NULL,
  `createDate` bigint(20) NOT NULL,
  `assign` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE IF NOT EXISTS `users` (
  `uuid` int(11) NOT NULL,
  `displayName` varchar(256) COLLATE utf8_bin NOT NULL,
  `name` varchar(128) COLLATE utf8_bin NOT NULL,
  `password` char(128) COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE IF NOT EXISTS `users_in_projects` (
  `user_id` int(11) NOT NULL,
  `project_id` int(11) NOT NULL,
  `enable` tinyint(1) NOT NULL,
  `expiration` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

ALTER TABLE `authority_role`
  ADD KEY `user_id` (`user_id`);

ALTER TABLE `projects`
  ADD PRIMARY KEY (`uuid`);

ALTER TABLE `todo_list`
  ADD PRIMARY KEY (`createDate`),
  ADD KEY `creator` (`creator`),
  ADD KEY `project` (`project`),
  ADD KEY `todo_list_users_uuid_fk` (`assign`);

ALTER TABLE `users`
  ADD PRIMARY KEY (`uuid`);

ALTER TABLE `users_in_projects`
  ADD KEY `project_id` (`project_id`),
  ADD KEY `user_id` (`user_id`);

ALTER TABLE `authority_role`
  ADD CONSTRAINT `authority_role_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`uuid`);

ALTER TABLE `todo_list`
  ADD CONSTRAINT `todo_list_ibfk_1` FOREIGN KEY (`creator`) REFERENCES `users` (`uuid`),
  ADD CONSTRAINT `todo_list_ibfk_2` FOREIGN KEY (`project`) REFERENCES `projects` (`uuid`),
  ADD CONSTRAINT `todo_list_users_uuid_fk` FOREIGN KEY (`assign`) REFERENCES `users` (`uuid`);

ALTER TABLE `users_in_projects`
  ADD CONSTRAINT `users_in_projects_ibfk_2` FOREIGN KEY (`project_id`) REFERENCES `projects` (`uuid`),
  ADD CONSTRAINT `users_in_projects_ibfk_3` FOREIGN KEY (`user_id`) REFERENCES `users` (`uuid`);
COMMIT;
