CREATE TABLE `etags` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`etag_id` VARCHAR(512) NOT NULL,
	`etag` VARCHAR(255) NOT NULL,
	`cached_until` DATETIME NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `etags_etag_id_unqiue_idx` (`etag_id`)
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;