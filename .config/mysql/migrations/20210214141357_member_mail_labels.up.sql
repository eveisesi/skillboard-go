CREATE TABLE `member_mail_labels` (
	`member_id` INT UNSIGNED NOT NULL,
	`labels` JSON NULL DEFAULT NULL,
	`total_unread_count` INT UNSIGNED NOT NULL DEFAULT '0',
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`member_id`) USING BTREE,
	CONSTRAINT `member_mail_labels_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;