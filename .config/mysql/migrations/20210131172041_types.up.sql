CREATE TABLE `types` (
	`id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`group_id` INT UNSIGNED NOT NULL,
	`published` TINYINT UNSIGNED NOT NULL DEFAULT '0',
	`capacity` FLOAT NOT NULL DEFAULT '0.00',
	`market_group_id` INT UNSIGNED NULL DEFAULT NULL,
	`mass` FLOAT NOT NULL DEFAULT '0.00',
	`packaged_volume` FLOAT NOT NULL DEFAULT '0.00',
	`portion_size` INT NULL DEFAULT NULL,
	`radius` FLOAT NULL DEFAULT NULL,
	`volume` FLOAT NOT NULL DEFAULT '0.00',
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `types_group_id_idx` (`group_id`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;