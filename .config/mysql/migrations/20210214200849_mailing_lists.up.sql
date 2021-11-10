CREATE TABLE `mailing_lists` (
	`mailing_list_id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(64) NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`mailing_list_id`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;