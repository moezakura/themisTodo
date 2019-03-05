START TRANSACTION;
ALTER TABLE `todo_list` DROP `name`;
ALTER TABLE `todo_list` DROP `status`;
ALTER TABLE `todo_list` DROP `deadline`;
ALTER TABLE `todo_list` DROP `description`;
ALTER TABLE `todo_list` DROP `assign`;
COMMIT;