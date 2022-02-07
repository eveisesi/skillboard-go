ALTER TABLE
    `users` DROP INDEX `character_id`,
ADD
    UNIQUE INDEX `character_id` (`character_id`) USING BTREE;