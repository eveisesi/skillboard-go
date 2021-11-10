CREATE TABLE `member_wallet_balance` (
    `member_id` INT UNSIGNED NOT NULL,
    `balance` FLOAT NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`) USING BTREE,
    CONSTRAINT `member_wallet_balance_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;