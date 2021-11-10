CREATE TABLE `member_home_clone` (
    `member_id` INT UNSIGNED NOT NULL,
    `location_id` BIGINT UNSIGNED NOT NULL,
    `location_type` ENUM('station', 'structure') NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`) USING BTREE,
    CONSTRAINT `member_home_clone_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;