CREATE TABLE `member_implants` (
    `member_id` INT UNSIGNED NOT NULL,
    `implant_id` INT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`, `implant_id`) USING BTREE,
    CONSTRAINT `member_implants_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT `member_implants_implant_id_foreign` FOREIGN KEY (`implant_id`) REFERENCES `athena`.`types` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;