CREATE TABLE `member_mailing_lists` (
	`member_id` INT UNSIGNED NOT NULL,
	`mailing_list_id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`member_id`, `mailing_list_id`) USING BTREE,
	CONSTRAINT `member_mailing_lists_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;