CREATE TABLE `member_contacts` (
    `member_id` INT UNSIGNED NOT NULL,
    `contact_id` INT UNSIGNED NOT NULL,
    `source_page` TINYINT UNSIGNED NOT NULL DEFAULT '1',
    `contact_type` VARCHAR(64) NOT NULL,
    `standing` FLOAT NOT NULL DEFAULT '0.00',
    `label_ids` JSON NOT NULL,
    `is_blocked` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `is_watched` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`, `contact_id`, `source_page`) USING BTREE,
    INDEX `member_contacts_source_page_idx` (`source_page`) USING BTREE,
    INDEX `member_contacts_contact_id_contact_type_idx` (`contact_id`, `contact_type`) USING BTREE,
    CONSTRAINT `member_contacts_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
)