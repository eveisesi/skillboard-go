CREATE TABLE `member_contact_labels` (
    `member_id` INT UNSIGNED NOT NULL,
    `label_id` BIGINT UNSIGNED NOT NULL,
    `label_name` VARCHAR(64) NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`, `label_id`) USING BTREE,
    CONSTRAINT `member_contact_labels_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;