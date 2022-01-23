CREATE TABLE `character_implants` (
    `character_id` BIGINT(20) UNSIGNED NOT NULL,
    `implant_id` INT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL,
    PRIMARY KEY (`character_id`, `implant_id`) USING BTREE,
    CONSTRAINT `character_implants_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `skillz`.`users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT `character_implants_implant_id_foreign` FOREIGN KEY (`implant_id`) REFERENCES `skillz`.`types` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;