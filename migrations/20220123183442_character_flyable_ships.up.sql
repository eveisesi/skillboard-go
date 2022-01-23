CREATE TABLE IF NOT EXISTS `character_flyable_ships` (
	`character_id` bigint(20) UNSIGNED NOT NULL,
	`ship_type_id` int(10) UNSIGNED NOT NULL,
	`flyable` tinyint(1) UNSIGNED NOT NULL DEFAULT '0',
	`created_at` datetime NOT NULL,
	PRIMARY KEY (`character_id`, `ship_type_id`),
	KEY `character_flyable_ships_ship_type_id` (`ship_type_id`),
	CONSTRAINT `character_flyable_ships_character_id` FOREIGN KEY (`character_id`) REFERENCES `users` (`character_id`) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT `character_flyable_ships_ship_type_id` FOREIGN KEY (`ship_type_id`) REFERENCES `types` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;