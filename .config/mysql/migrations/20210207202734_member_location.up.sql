CREATE TABLE `member_location` (
    `member_id` INT UNSIGNED NOT NULL,
    `solar_system_id` INT UNSIGNED NOT NULL,
    `station_id` INT UNSIGNED NULL DEFAULT NULL,
    `structure_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`) USING BTREE,
    CONSTRAINT `member_location_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;