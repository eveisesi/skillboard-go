CREATE TABLE `character_skills` (
    `character_id` BIGINT(20) UNSIGNED NOT NULL,
    `skill_id` INT UNSIGNED NOT NULL,
    `active_skill_level` TINYINT UNSIGNED NOT NULL,
    `skillpoints_in_skill` INT UNSIGNED NOT NULL,
    `trained_skill_level` TINYINT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`character_id`, `skill_id`) USING BTREE,
    CONSTRAINT `character_skills_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `skillz`.`users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;