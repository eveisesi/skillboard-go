CREATE TABLE `member_online` (
    `member_id` INT UNSIGNED NOT NULL,
    `last_login` TIMESTAMP NULL DEFAULT NULL,
    `last_logout` TIMESTAMP NULL DEFAULT NULL,
    `logins` INT UNSIGNED NOT NULL,
    `online` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`) USING BTREE,
    CONSTRAINT `member_online_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;