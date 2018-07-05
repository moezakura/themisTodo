ALTER TABLE `users`
    ADD COLUMN `password_version` INT NOT NULL DEFAULT 0;