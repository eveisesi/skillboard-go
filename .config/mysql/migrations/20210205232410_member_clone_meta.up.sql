CREATE TABLE `member_clone_meta` (
    `member_id` INT UNSIGNED NOT NULL,
    `last_clone_jump_date` TIMESTAMP NULL DEFAULT NULL,
    `last_station_change_date` TIMESTAMP NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`) USING BTREE,
    CONSTRAINT `member_clone_meta_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;