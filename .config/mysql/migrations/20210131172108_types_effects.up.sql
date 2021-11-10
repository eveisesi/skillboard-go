CREATE TABLE `type_effects` (
    `type_id` INT UNSIGNED NOT NULL,
    `effect_id` INT UNSIGNED NOT NULL,
    `is_default` TINYINT NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`type_id`, `effect_id`) USING BTREE,
    CONSTRAINT `type_effects_type_id_types_id_foreign` FOREIGN KEY (`type_id`) REFERENCES `athena`.`types` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;