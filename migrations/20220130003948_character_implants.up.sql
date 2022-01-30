CREATE TABLE `character_implants` (
    `character_id` bigint(20) UNSIGNED NOT NULL,
    `implant_id` int(10) UNSIGNED NOT NULL,
    `slot` int(10) UNSIGNED NOT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`character_id`, `implant_id`),
    KEY `character_implants_implant_id` (`implant_id`),
    CONSTRAINT `character_implants_implant_id_foreign` FOREIGN KEY (`implant_id`) REFERENCES `types` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT `character_implants_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `users` (`character_id`) ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;