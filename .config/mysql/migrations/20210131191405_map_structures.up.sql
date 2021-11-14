CREATE TABLE `map_structures` (
    `id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `solar_system_id` INT UNSIGNED NOT NULL,
    `type_id` INT UNSIGNED NOT NULL,
    `owner_id` INT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `structures_solar_system_id_idx` (`solar_system_id`) USING BTREE,
    INDEX `structures_type_id_idx` (`type_id`) USING BTREE,
    INDEX `structures_owner_id_idx` (`owner_id`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;