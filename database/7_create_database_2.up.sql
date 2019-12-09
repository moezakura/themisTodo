-- phpMyAdmin SQL Dump
-- version 4.8.3
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: 2019 年 12 月 09 日 06:17
-- サーバのバージョン： 5.5.64-MariaDB
-- PHP Version: 7.3.0RC1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT = @@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS = @@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION = @@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `themis_todo`
--

-- --------------------------------------------------------

--
-- テーブルの構造 `authority_role`
--

CREATE TABLE IF NOT EXISTS `authority_role`
(
    `user_id` int(11)    NOT NULL,
    `type`    int(11)    NOT NULL,
    `enable`  tinyint(1) NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_bin;

-- --------------------------------------------------------

--
-- テーブルの構造 `projects`
--

CREATE TABLE IF NOT EXISTS `projects`
(
    `uuid`        int(11)                        NOT NULL,
    `name`        varchar(256) COLLATE utf8_bin  NOT NULL,
    `description` varchar(1024) COLLATE utf8_bin NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_bin;

--
-- テーブルの構造 `todo_list`
--

CREATE TABLE IF NOT EXISTS `todo_list`
(
    `id`         int(11)    NOT NULL,
    `project`    int(11)    NOT NULL,
    `creator`    int(11)    NOT NULL,
    `createDate` bigint(20) NOT NULL,
    `adopted`    bigint(20) DEFAULT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_bin;

-- --------------------------------------------------------

--
-- テーブルの構造 `todo_list_history`
--

CREATE TABLE IF NOT EXISTS `todo_list_history`
(
    `name`        varchar(1024) COLLATE utf8_bin NOT NULL,
    `editor`      int(11)                        NOT NULL,
    `status`      int(11)                        NOT NULL,
    `deadline`    date                           NOT NULL,
    `description` text COLLATE utf8_bin          NOT NULL,
    `createDate`  bigint(20)                     NOT NULL,
    `updateDate`  bigint(20)                     NOT NULL,
    `assign`      int(11)                        NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_bin;

-- --------------------------------------------------------

--
-- テーブルの構造 `todo_timer`
--

CREATE TABLE IF NOT EXISTS `todo_timer`
(
    `id`         int(11)    NOT NULL,
    `createDate` bigint(20) NOT NULL,
    `assign`     int(11)    NOT NULL,
    `startDate`  datetime   NOT NULL,
    `endDate`    datetime   NOT NULL,
    `note`       blob
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_bin;

-- --------------------------------------------------------

--
-- テーブルの構造 `users`
--

CREATE TABLE IF NOT EXISTS `users`
(
    `uuid`        int(11)                       NOT NULL,
    `displayName` varchar(256) COLLATE utf8_bin NOT NULL,
    `name`        varchar(128) COLLATE utf8_bin NOT NULL,
    `icon_path`   varchar(48) COLLATE utf8_bin  NOT NULL DEFAULT 'noIcon',
    `password`    char(128) COLLATE utf8_bin    NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_bin;

-- --------------------------------------------------------

--
-- テーブルの構造 `users_in_projects`
--

CREATE TABLE IF NOT EXISTS `users_in_projects`
(
    `user_id`    int(11)    NOT NULL,
    `project_id` int(11)    NOT NULL,
    `enable`     tinyint(1) NOT NULL,
    `expiration` datetime DEFAULT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_bin;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `authority_role`
--
ALTER TABLE `authority_role`
    ADD KEY `user_id` (`user_id`);

--
-- Indexes for table `projects`
--
ALTER TABLE `projects`
    ADD PRIMARY KEY (`uuid`);

--
-- Indexes for table `todo_list`
--
ALTER TABLE `todo_list`
    ADD PRIMARY KEY (`createDate`),
    ADD KEY `creator` (`creator`),
    ADD KEY `project` (`project`),
    ADD KEY `todo_list_ibfk_3` (`adopted`);

--
-- Indexes for table `todo_list_history`
--
ALTER TABLE `todo_list_history`
    ADD PRIMARY KEY (`updateDate`),
    ADD KEY `creator` (`editor`),
    ADD KEY `todo_list_history_users_uuid_fk` (`assign`),
    ADD KEY `todo_list_history_createDate_fk` (`createDate`);

--
-- Indexes for table `todo_timer`
--
ALTER TABLE `todo_timer`
    ADD PRIMARY KEY (`id`),
    ADD KEY `todo_timer_users_uuid_fk` (`assign`),
    ADD KEY `todo_timer_createDate_fk` (`createDate`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
    ADD PRIMARY KEY (`uuid`);

--
-- Indexes for table `users_in_projects`
--
ALTER TABLE `users_in_projects`
    ADD KEY `project_id` (`project_id`),
    ADD KEY `user_id` (`user_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `projects`
--
ALTER TABLE `projects`
    MODIFY `uuid` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `todo_timer`
--
ALTER TABLE `todo_timer`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
    MODIFY `uuid` int(11) NOT NULL AUTO_INCREMENT;

--
-- ダンプしたテーブルの制約
--

--
-- テーブルの制約 `authority_role`
--
ALTER TABLE `authority_role`
    ADD CONSTRAINT `authority_role_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`uuid`);

--
-- テーブルの制約 `todo_list`
--
ALTER TABLE `todo_list`
    ADD CONSTRAINT `todo_list_ibfk_1` FOREIGN KEY (`creator`) REFERENCES `users` (`uuid`),
    ADD CONSTRAINT `todo_list_ibfk_2` FOREIGN KEY (`project`) REFERENCES `projects` (`uuid`),
    ADD CONSTRAINT `todo_list_ibfk_3` FOREIGN KEY (`adopted`) REFERENCES `todo_list_history` (`updateDate`) ON UPDATE SET NULL;

--
-- テーブルの制約 `todo_list_history`
--
ALTER TABLE `todo_list_history`
    ADD CONSTRAINT `todo_list_history_createDate_fk` FOREIGN KEY (`createDate`) REFERENCES `todo_list` (`createDate`),
    ADD CONSTRAINT `todo_list_history_ibfk_1` FOREIGN KEY (`editor`) REFERENCES `users` (`uuid`),
    ADD CONSTRAINT `todo_list_history_users_uuid_fk` FOREIGN KEY (`assign`) REFERENCES `users` (`uuid`);

--
-- テーブルの制約 `todo_timer`
--
ALTER TABLE `todo_timer`
    ADD CONSTRAINT `todo_timer_createDate_fk` FOREIGN KEY (`createDate`) REFERENCES `todo_list` (`createDate`),
    ADD CONSTRAINT `todo_timer_users_uuid_fk` FOREIGN KEY (`assign`) REFERENCES `users` (`uuid`);

--
-- テーブルの制約 `users_in_projects`
--
ALTER TABLE `users_in_projects`
    ADD CONSTRAINT `users_in_projects_ibfk_2` FOREIGN KEY (`project_id`) REFERENCES `projects` (`uuid`),
    ADD CONSTRAINT `users_in_projects_ibfk_3` FOREIGN KEY (`user_id`) REFERENCES `users` (`uuid`);
COMMIT;
