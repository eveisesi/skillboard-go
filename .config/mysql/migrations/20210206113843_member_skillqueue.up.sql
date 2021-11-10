CREATE TABLE `member_skillqueue` (
    `member_id` INT UNSIGNED NOT NULL,
    `queue_position` TINYINT UNSIGNED NOT NULL,
    `skill_id` INT UNSIGNED NOT NULL,
    `start_date` TIMESTAMP NULL DEFAULT NULL,
    `finish_date` TIMESTAMP NULL DEFAULT NULL,
    `finished_level` TINYINT UNSIGNED NOT NULL,
    `training_start_sp` INT UNSIGNED NULL DEFAULT NULL,
    `level_start_sp` INT UNSIGNED NULL DEFAULT NULL,
    `level_end_sp` INT UNSIGNED NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`, `queue_position`) USING BTREE,
    INDEX `member_skillqueue_skill_id_idx` (`skill_id`),
    INDEX `member_skillqueue_start_date_idx` (`start_date`),
    INDEX `member_skillqueue_finish_date_idx` (`finish_date`),
    CONSTRAINT `member_skillqueue_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;