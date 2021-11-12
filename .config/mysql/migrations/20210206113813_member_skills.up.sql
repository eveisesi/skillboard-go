CREATE TABLE `member_skills` (
    `character_id` BIGINT(20) UNSIGNED NOT NULL,
    `skill_id` INT UNSIGNED NOT NULL,
    `active_skill_level` TINYINT UNSIGNED NOT NULL,
    `skillpoints_in_skill` INT UNSIGNED NOT NULL,
    `trained_skill_level` TINYINT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`member_id`, `skill_id`) USING BTREE,
    CONSTRAINT `member_skills_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `skillz`.`users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;