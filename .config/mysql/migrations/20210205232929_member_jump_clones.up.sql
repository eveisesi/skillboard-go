CREATE TABLE `member_jump_clones` (
    `member_id` INT UNSIGNED NOT NULL,
    `jump_clone_id` INT UNSIGNED NOT NULL,
    `location_id` BIGINT UNSIGNED NOT NULL,
    `location_type` ENUM('station', 'structure') NOT NULL,
    `implants` JSON NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`, `jump_clone_id`) USING BTREE,
    CONSTRAINT `member_jump_clones_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;