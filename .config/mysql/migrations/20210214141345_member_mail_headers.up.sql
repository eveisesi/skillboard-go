CREATE TABLE `member_mail_headers` (
	`member_id` INT UNSIGNED NOT NULL,
	`mail_id` INT UNSIGNED NOT NULL,
	`labels` JSON NOT NULL,
	`is_read` TINYINT NOT NULL DEFAULT '0',
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`member_id`, `mail_id`) USING BTREE,
	CONSTRAINT `member_mail_headers_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT `member_mail_headers_mail_id_foreign` FOREIGN KEY (`mail_id`) REFERENCES `athena`.`mail_headers` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;