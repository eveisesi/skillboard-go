CREATE TABLE `character_home_clone` (
    `character_id` BIGINT(20) UNSIGNED NOT NULL,
    `location_id` BIGINT UNSIGNED NOT NULL,
    `location_type` ENUM('station', 'structure') NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`character_id`) USING BTREE,
    CONSTRAINT `character_home_clone_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `skillz`.`users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;