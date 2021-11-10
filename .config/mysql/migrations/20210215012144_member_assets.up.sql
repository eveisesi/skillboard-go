CREATE TABLE `member_assets` (
	`member_id` INT UNSIGNED NOT NULL,
	`item_id` BIGINT UNSIGNED NOT NULL,
	`type_id` INT UNSIGNED NOT NULL,
	`location_id` BIGINT UNSIGNED NOT NULL,
	`location_type` VARCHAR(64) NOT NULL,
	`location_flag` VARCHAR(64) NOT NULL,
	`quantity` INT NOT NULL,
	`is_blueprint_copy` TINYINT UNSIGNED NOT NULL DEFAULT '0',
	`is_singleton` TINYINT UNSIGNED NOT NULL DEFAULT '0',
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`member_id`, `item_id`) USING BTREE,
	INDEX `member_assets_type_id_idx` (`type_id`) USING BTREE,
	INDEX `member_assets_location_id_idx` (`location_id`) USING BTREE,
	INDEX `member_assets_location_type_idx` (`location_type`) USING BTREE,
	INDEX `member_assets_location_flag_idx` (`location_flag`) USING BTREE,
	CONSTRAINT `member_assets_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;