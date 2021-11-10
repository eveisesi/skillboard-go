CREATE TABLE `member_contract_items` (
    `member_id` INT UNSIGNED NOT NULL,
    `contract_id` INT UNSIGNED NOT NULL,
    `record_id` INT UNSIGNED NOT NULL,
    `type_id` INT UNSIGNED NOT NULL,
    `quantity` INT UNSIGNED NOT NULL,
    `raw_quantity` INT NOT NULL,
    `is_included` TINYINT UNSIGNED NOT NULL,
    `is_singleton` TINYINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`, `contract_id`, `record_id`) USING BTREE,
    INDEX `member_contract_items_type_id_idx` (`type_id`) USING BTREE,
    CONSTRAINT `member_contract_items_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;