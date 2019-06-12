CREATE TABLE `todo_timer`
(
    `id`         INT      NOT NULL AUTO_INCREMENT,
    `createDate` BIGINT   NOT NULL,
    `assign`     INT      NOT NULL,
    `startDate`  DATETIME NOT NULL,
    `endDate`    DATETIME NOT NULL,
    `note`       BLOB     NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

ALTER TABLE `todo_timer`
    ADD CONSTRAINT `todo_timer_users_uuid_fk` FOREIGN KEY (`assign`) REFERENCES `users` (`uuid`),
    ADD CONSTRAINT `todo_timer_createDate_fk` FOREIGN KEY (`createDate`) REFERENCES `todo_list` (`createDate`);