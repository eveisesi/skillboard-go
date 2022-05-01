CREATE TABLE `ship_flight_requirements` (
    `group_id` int unsigned not null,
    `ship_id` int unsigned not null,
    `skill_id` int unsigned not null,
    `minimum_skill_level` int unsigned not null,
    `created_at` timestamp not null,
    PRIMARY KEY (`ship_id`, `skill_id`),
    CONSTRAINT `ship_flight_requirements_ship_id_foreign` FOREIGN KEY (`ship_id`) REFERENCES `types` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT `ship_flight_requirements_skill_id_foreign` FOREIGN KEY (`skill_id`) REFERENCES `types` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT `ship_flight_requirements_group_id_foreign` FOREIGN KEY (`group_id`) REFERENCES `type_groups` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;