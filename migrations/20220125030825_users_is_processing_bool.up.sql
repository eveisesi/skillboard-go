ALTER TABLE
    users
ADD
    COLUMN `is_processing` TINYINT(1) NOT NULL DEFAULT '0'
AFTER
    `is_new`;