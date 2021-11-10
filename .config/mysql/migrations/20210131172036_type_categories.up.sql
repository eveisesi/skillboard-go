CREATE TABLE `type_categories` (
	`id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`published` TINYINT NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;