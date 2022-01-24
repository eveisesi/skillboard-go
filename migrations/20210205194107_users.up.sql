CREATE TABLE `users` (
    `id` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_unicode_ci',
    `character_id` BIGINT UNSIGNED NOT NULL,
    `access_token` TEXT NOT NULL COLLATE 'utf8mb4_unicode_ci',
    `refresh_token` TEXT NOT NULL COLLATE 'utf8mb4_unicode_ci',
    `expires` DATETIME NOT NULL,
    `owner_hash` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_unicode_ci',
    `scopes` JSON NOT NULL,
    `is_new` TINYINT(1) NOT NULL DEFAULT '1',
    `disabled` TINYINT(3) UNSIGNED NOT NULL DEFAULT '0',
    `disabled_reason` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_unicode_ci',
    `disabled_timestamp` DATETIME NULL DEFAULT NULL,
    `last_login` DATETIME NOT NULL,
    `last_processed` DATETIME NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `character_id` (`character_id`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;