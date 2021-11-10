CREATE TABLE `member_ship`(
    `member_id` INT UNSIGNED NOT NULL,
    `ship_item_id` BIGINT UNSIGNED NOT NULL,
    `ship_type_id` INT UNSIGNED NOT NULL,
    `ship_name` VARCHAR(64) NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`) USING BTREE,
    CONSTRAINT `member_ship_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;