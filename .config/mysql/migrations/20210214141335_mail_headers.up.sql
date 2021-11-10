CREATE TABLE `mail_headers` (
	`id` INT UNSIGNED NOT NULL,
	`sender_id` INT UNSIGNED NOT NULL,
	`sender_type` VARCHAR(32) NULL DEFAULT NULL,
	`subject` VARCHAR(255) NULL DEFAULT NULL,
	`body` TEXT NULL DEFAULT NULL,
	`sent` TIMESTAMP NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `mail_headers_sender_id_idx` (`sender_id`) USING BTREE,
	INDEX `mail_headers_sender_type_idx` (`sender_type`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;