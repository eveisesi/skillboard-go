CREATE TABLE `members` (
    `id` INT UNSIGNED NOT NULL,
    `main_id` INT UNSIGNED NULL DEFAULT NULL,
    `access_token` TEXT NULL DEFAULT NULL,
    `refresh_token` VARCHAR(128) NULL DEFAULT NULL,
    `expires` TIMESTAMP NULL DEFAULT NULL,
    `owner_hash` VARCHAR(128) NULL DEFAULT NULL,
    `scopes` JSON NOT NULL,
    `disabled` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `disabled_reason` VARCHAR(255) NULL DEFAULT NULL,
    `disabled_timestamp` TIMESTAMP NULL DEFAULT NULL,
    `last_login` TIMESTAMP NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY(`id`) USING BTREE,
    INDEX `members_main_id_idx` (`main_id`)
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;