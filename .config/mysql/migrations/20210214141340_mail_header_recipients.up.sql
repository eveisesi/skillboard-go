CREATE TABLE `mail_header_recipients` (
	`mail_id` INT UNSIGNED NOT NULL,
	`recipient_id` INT UNSIGNED NOT NULL,
	`recipient_type` VARCHAR(32) NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`mail_id`, `recipient_id`) USING BTREE,
	CONSTRAINT `mail_header_mail_id_foreign` FOREIGN KEY (`mail_id`) REFERENCES `athena`.`mail_headers` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;