CREATE TABLE `member_contact_labels` (
    `character_id` BIGINT(20) UNSIGNED NOT NULL,
    `label_id` BIGINT UNSIGNED NOT NULL,
    `label_name` VARCHAR(64) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`member_id`, `label_id`) USING BTREE,
    CONSTRAINT `member_contact_labels_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `skillz`.`users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;