CREATE TABLE `map_solar_systems` (
    `id` INT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `constellation_id` INT UNSIGNED NOT NULL,
    `security_class` VARCHAR(255) NULL DEFAULT NULL,
    `security_status` FLOAT NOT NULL,
    `star_id` INT UNSIGNED NULL DEFAULT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `solar_systems_constellation_id_idx` (`constellation_id`) USING BTREE,
    INDEX `solar_systems_star_id_idx` (`star_id`) USING BTREE,
    INDEX `solar_systems_security_status` (`security_status`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;