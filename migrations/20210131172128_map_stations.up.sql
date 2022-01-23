CREATE TABLE `map_stations` (
    `id` INT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `system_id` INT UNSIGNED NOT NULL,
    `type_id` INT UNSIGNED NOT NULL,
    `race_id` INT UNSIGNED NULL DEFAULT NULL,
    `owner_corporation_id` INT UNSIGNED NULL DEFAULT NULL,
    `max_dockable_ship_volume` FLOAT NOT NULL,
    `office_rental_cost` FLOAT NOT NULL,
    `reprocessing_efficiency` FLOAT NOT NULL,
    `reprocessing_stations_take` FLOAT NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `stations_system_id_idx` (`system_id`) USING BTREE,
    INDEX `stations_type_id_idx` (`type_id`) USING BTREE,
    INDEX `stations_race_id_idx` (`race_id`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;