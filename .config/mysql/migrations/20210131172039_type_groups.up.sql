CREATE TABLE `type_groups` (
	`id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`published` TINYINT NOT NULL,
	`category_id` INT UNSIGNED NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `type_groups_category_id_idx` (`category_id`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;