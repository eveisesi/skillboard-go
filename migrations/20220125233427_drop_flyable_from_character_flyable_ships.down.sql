ALTER TABLE
    `character_flyable_ships`
ADD
    COLUMN `flyable` tinyint(1) UNSIGNED NOT NULL DEFAULT '0'
AFTER
    `ship_type_id`;