CREATE TABLE `character_clone_meta` (
    `character_id` BIGINT(20) UNSIGNED NOT NULL,
    `last_clone_jump_date` TIMESTAMP NULL DEFAULT NULL,
    `last_station_change_date` TIMESTAMP NULL DEFAULT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`character_id`) USING BTREE,
    CONSTRAINT `character_clone_meta_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `skillz`.`users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;