CREATE TABLE `bloodlines` (
	`id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`race_id` INT UNSIGNED NOT NULL,
	`corporation_id` INT UNSIGNED NOT NULL,
	`ship_type_id` INT UNSIGNED NOT NULL,
	`charisma` INT UNSIGNED NOT NULL,
	`intelligence` INT UNSIGNED NOT NULL,
	`memory` INT UNSIGNED NOT NULL,
	`perception` INT UNSIGNED NOT NULL,
	`willpower` INT UNSIGNED NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `bloodlines_race_id_idx` (`race_id`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;