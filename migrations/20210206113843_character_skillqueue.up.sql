CREATE TABLE `character_skillqueue` (
    `character_id` BIGINT(20) UNSIGNED NOT NULL,
    `queue_position` TINYINT UNSIGNED NOT NULL,
    `skill_id` INT UNSIGNED NOT NULL,
    `start_date` TIMESTAMP NULL DEFAULT NULL,
    `finish_date` TIMESTAMP NULL DEFAULT NULL,
    `finished_level` TINYINT UNSIGNED NOT NULL,
    `training_start_sp` INT UNSIGNED NULL DEFAULT NULL,
    `level_start_sp` INT UNSIGNED NULL DEFAULT NULL,
    `level_end_sp` INT UNSIGNED NULL DEFAULT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`character_id`, `queue_position`) USING BTREE,
    INDEX `character_skillqueue_skill_id_idx` (`skill_id`),
    INDEX `character_skillqueue_start_date_idx` (`start_date`),
    INDEX `character_skillqueue_finish_date_idx` (`finish_date`),
    CONSTRAINT `character_skillqueue_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `skillz`.`users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;