CREATE TABLE `character_contacts` (
    `character_id` BIGINT(20) UNSIGNED NOT NULL,
    `contact_id` INT UNSIGNED NOT NULL,
    `source_page` TINYINT UNSIGNED NOT NULL DEFAULT '1',
    `contact_type` VARCHAR(64) NOT NULL,
    `standing` FLOAT NOT NULL DEFAULT '0.00',
    `label_ids` JSON NOT NULL,
    `is_blocked` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `is_watched` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`character_id`, `contact_id`, `source_page`) USING BTREE,
    INDEX `character_contacts_source_page_idx` (`source_page`) USING BTREE,
    INDEX `character_contacts_contact_id_contact_type_idx` (`contact_id`, `contact_type`) USING BTREE,
    CONSTRAINT `character_contacts_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `skillz`.`users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE
)