CREATE TABLE `member_attributes` (
    `member_id` INT UNSIGNED NOT NULL,
    `charisma` TINYINT UNSIGNED NOT NULL,
    `intelligence` TINYINT UNSIGNED NOT NULL,
    `memory` TINYINT UNSIGNED NOT NULL,
    `perception` TINYINT UNSIGNED NOT NULL,
    `willpower` TINYINT UNSIGNED NOT NULL,
    `bonus_remaps` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `accured_remap_cooldown_date` TIMESTAMP NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`) USING BTREE,
    CONSTRAINT `member_skill_attributes_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;