CREATE TABLE `todo_list_history` (
  `id` int(11) NOT NULL,
  `project` int(11) NOT NULL,
  `name` varchar(1024) COLLATE utf8_bin NOT NULL,
  `creator` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  `deadline` date NOT NULL,
  `description` text COLLATE utf8_bin NOT NULL,
  `createDate` bigint(20) NOT NULL,
  `parent_createDate` bigint(20) NOT NULL,
  `assign` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

ALTER TABLE `todo_list_history`
  ADD PRIMARY KEY (`createDate`),
  ADD KEY `creator` (`creator`),
  ADD KEY `project` (`project`),
  ADD KEY `todo_list_history_users_uuid_fk` (`assign`),
  ADD KEY `todo_list_history_parent_createDate_fk` (`parent_createDate`);

ALTER TABLE `todo_list_history`
  ADD CONSTRAINT `todo_list_history_ibfk_1` FOREIGN KEY (`creator`) REFERENCES `users` (`uuid`),
  ADD CONSTRAINT `todo_list_history_ibfk_2` FOREIGN KEY (`project`) REFERENCES `projects` (`uuid`),
  ADD CONSTRAINT `todo_list_history_users_uuid_fk` FOREIGN KEY (`assign`) REFERENCES `users` (`uuid`),
  ADD CONSTRAINT `todo_list_history_parent_createDate_fk` FOREIGN KEY (`parent_createDate`) REFERENCES `todo_list` (`createDate`);
