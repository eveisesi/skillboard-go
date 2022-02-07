ALTER TABLE
    `users` CHANGE COLUMN `owner_hash` `owner_hash` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_unicode_ci'
AFTER
    `expires`;