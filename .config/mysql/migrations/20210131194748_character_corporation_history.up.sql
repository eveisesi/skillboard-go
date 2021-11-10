CREATE TABLE `character_corporation_history` (
    `id` INT UNSIGNED NOT NULL,
    `record_id` INT UNSIGNED NOT NULL,
    `corporation_id` INT UNSIGNED NOT NULL,
    `is_deleted` TINYINT NOT NULL DEFAULT '0',
    `start_date` TIMESTAMP NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`, `record_id`) USING BTREE,
    CONSTRAINT `character_corporation_history_id_foreign` FOREIGN KEY (`id`) REFERENCES `athena`.`characters` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;