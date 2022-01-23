CREATE TABLE `alliances` (
    `id` INT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `ticker` VARCHAR(255) NOT NULL,
    `creator_id` INT UNSIGNED NOT NULL,
    `creator_corporation_id` INT UNSIGNED NOT NULL,
    `executor_corporation_id` INT UNSIGNED NULL DEFAULT NULL,
    `is_closed` TINYINT NOT NULL DEFAULT '0',
    `date_founded` TIMESTAMP NULL DEFAULT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `alliances_creator_id` (`creator_id`) USING BTREE,
    INDEX `alliances_creator_corporation_id` (`creator_corporation_id`) USING BTREE,
    INDEX `alliances_executor_corporation_id` (`executor_corporation_id`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;