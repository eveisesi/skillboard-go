CREATE TABLE `character_skill_meta` (
    `character_id` BIGINT(20) UNSIGNED NOT NULL,
    `total_sp` INT UNSIGNED NOT NULL,
    `unallocated_sp` INT UNSIGNED NULL DEFAULT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`character_id`) USING BTREE,
    CONSTRAINT `character_skill_meta_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `skillboard`.`users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;