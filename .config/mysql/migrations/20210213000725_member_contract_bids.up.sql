CREATE TABLE `member_contract_bids` (
	`member_id` INT UNSIGNED NOT NULL,
	`contract_id` INT UNSIGNED NOT NULL,
	`bid_id` INT UNSIGNED NOT NULL,
	`bidder_id` INT UNSIGNED NOT NULL,
	`amount` FLOAT NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`member_id`, `contract_id`, `bid_id`) USING BTREE,
	INDEX `member_contract_bids_bidder_id_idx` (`bidder_id`) USING BTREE,
	CONSTRAINT `member_contract_bids_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;