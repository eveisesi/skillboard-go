CREATE TABLE `member_skill_properties` (
    `member_id` INT UNSIGNED NOT NULL,
    `total_sp` INT UNSIGNED NOT NULL,
    `unallocated_sp` INT UNSIGNED NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`) USING BTREE,
    CONSTRAINT `member_skill_meta_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;