CREATE TABLE `character_attributes` (
    `character_id` BIGINT(20) UNSIGNED NOT NULL,
    `charisma` TINYINT UNSIGNED NOT NULL,
    `intelligence` TINYINT UNSIGNED NOT NULL,
    `memory` TINYINT UNSIGNED NOT NULL,
    `perception` TINYINT UNSIGNED NOT NULL,
    `willpower` TINYINT UNSIGNED NOT NULL,
    `bonus_remaps` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `last_remap_date` DATETIME NULL DEFAULT NULL,
    `accrued_remap_cooldown_date` DATETIME NULL DEFAULT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`character_id`) USING BTREE,
    CONSTRAINT `character_skill_attributes_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `skillboard`.`users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;