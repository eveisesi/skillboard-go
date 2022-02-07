ALTER TABLE
    `users`
ADD
    COLUMN `uuid` VARCHAR(128) NOT NULL
AFTER
    `id`;