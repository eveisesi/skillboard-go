ALTER TABLE
    `user_settings`
ADD
    COLUMN `hide_implants` tinyint(1) UNSIGNED NOT NULL DEFAULT '1'
AFTER
    `hide_attributes`;