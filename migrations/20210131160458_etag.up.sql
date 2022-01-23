CREATE TABLE `etags` (
	`path` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_unicode_ci',
	`etag` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_unicode_ci',
	`cached_until` DATETIME NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	INDEX `path` (`path`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;