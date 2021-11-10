CREATE TABLE `member_skills` (
    `member_id` INT UNSIGNED NOT NULL,
    `skill_id` INT UNSIGNED NOT NULL,
    `active_skill_level` TINYINT UNSIGNED NOT NULL,
    `skillpoints_in_skill` INT UNSIGNED NOT NULL,
    `trained_skill_level` TINYINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`, `skill_id`) USING BTREE,
    CONSTRAINT `member_skills_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;