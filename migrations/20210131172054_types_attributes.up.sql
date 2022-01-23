CREATE TABLE `type_attributes` (
    `type_id` INT UNSIGNED NOT NULL,
    `attribute_id` INT UNSIGNED NOT NULL,
    `value` FLOAT NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`type_id`, `attribute_id`) USING BTREE,
    CONSTRAINT `types_attributes_type_id_types_id_foreign` FOREIGN KEY (`type_id`) REFERENCES `skillboard`.`types` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;