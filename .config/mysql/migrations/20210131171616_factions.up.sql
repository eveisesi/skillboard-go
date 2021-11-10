CREATE TABLE `factions` (
	`id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`is_unique` TINYINT UNSIGNED NOT NULL DEFAULT '0',
	`size_factor` FLOAT NOT NULL,
	`station_count` INT UNSIGNED NOT NULL,
	`station_system_count` INT UNSIGNED NOT NULL,
	`corporation_id` INT UNSIGNED NULL DEFAULT NULL,
	`militia_corporation_id` INT UNSIGNED NULL DEFAULT NULL,
	`solar_system_id` INT UNSIGNED NULL DEFAULT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `factions_solar_system_id_idx` (`solar_system_id`),
	INDEX `factions_corporation_id_idx` (`corporation_id`),
	INDEX `factions_militia_corporation_id_idx` (`militia_corporation_id`)
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;