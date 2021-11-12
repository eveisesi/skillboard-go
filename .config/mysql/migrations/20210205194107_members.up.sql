CREATE TABLE `users` (
    `id` INT UNSIGNED NOT NULL,
    `character_id` INT UNSIGNED NOT NULL,
    `access_token` TEXT NOT NULL,
    `refresh_token` TEXT NOT NULL,
    `expires` DATETIME NOT NULL,
    `owner_hash` VARCHAR(128) NOT NULL,
    `scopes` JSON NOT NULL,
    `disabled` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `disabled_reason` VARCHAR(255) NULL DEFAULT NULL,
    `disabled_timestamp` DATETIME NULL DEFAULT NULL,
    `last_login` DATETIME NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY(`id`) USING BTREE,
    INDEX `users_character_id_idx` (`character_id`)
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;