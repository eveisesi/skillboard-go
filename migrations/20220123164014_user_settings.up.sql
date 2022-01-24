CREATE TABLE `user_settings` (
	`user_id` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_unicode_ci',
	`visibility` VARCHAR(32) NOT NULL COLLATE 'utf8mb4_unicode_ci',
	`visibility_token` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_unicode_ci',
	`hide_skills` TINYINT(1) NOT NULL DEFAULT '0',
	`hide_queue` TINYINT(1) NOT NULL DEFAULT '0',
	`hide_flyable` TINYINT(1) NOT NULL DEFAULT '0',
	`hide_attributes` TINYINT(1) NOT NULL DEFAULT '0',
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (`user_id`) USING BTREE,
	CONSTRAINT `user_settings_user_id_users_id` FOREIGN KEY (`user_id`) REFERENCES `skillboard`.`users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;