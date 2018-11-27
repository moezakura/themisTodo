CREATE TABLE `todo_list_history` (
  `name` varchar(1024) COLLATE utf8_bin NOT NULL,
  `editor` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  `deadline` date NOT NULL,
  `description` text COLLATE utf8_bin NOT NULL,
  `createDate` bigint(20) NOT NULL,
  `updateDate` bigint(20) NOT NULL,
  `assign` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

ALTER TABLE `todo_list_history`
  ADD PRIMARY KEY (`updateDate`),
  ADD KEY `creator` (`editor`),
  ADD KEY `todo_list_history_users_uuid_fk` (`assign`),
  ADD KEY `todo_list_history_createDate_fk` (`createDate`);

ALTER TABLE `todo_list_history`
  ADD CONSTRAINT `todo_list_history_ibfk_1` FOREIGN KEY (`editor`) REFERENCES `users` (`uuid`),
  ADD CONSTRAINT `todo_list_history_users_uuid_fk` FOREIGN KEY (`assign`) REFERENCES `users` (`uuid`),
  ADD CONSTRAINT `todo_list_history_createDate_fk` FOREIGN KEY (`createDate`) REFERENCES `todo_list` (`createDate`);
