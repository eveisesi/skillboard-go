-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               5.7.33 - MySQL Community Server (GPL)
-- Server OS:                    Linux
-- HeidiSQL Version:             11.3.0.6295
-- --------------------------------------------------------
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */
;

/*!40101 SET NAMES utf8 */
;

/*!50503 SET NAMES utf8mb4 */
;

/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */
;

/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */
;

/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */
;

-- Dumping database structure for skillz
CREATE DATABASE IF NOT EXISTS `skillz`
/*!40100 DEFAULT CHARACTER SET latin1 */
;

USE `skillz`;

-- Dumping structure for table skillz.alliances
CREATE TABLE IF NOT EXISTS `alliances` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `ticker` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `creator_id` int(10) UNSIGNED NOT NULL,
    `creator_corporation_id` int(10) UNSIGNED NOT NULL,
    `executor_corporation_id` int(10) UNSIGNED DEFAULT NULL,
    `is_closed` tinyint(4) NOT NULL DEFAULT '0',
    `date_founded` timestamp NULL DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `alliances_creator_id` (`creator_id`) USING BTREE,
    KEY `alliances_creator_corporation_id` (`creator_corporation_id`) USING BTREE,
    KEY `alliances_executor_corporation_id` (`executor_corporation_id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.bloodlines
CREATE TABLE IF NOT EXISTS `bloodlines` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `race_id` int(10) UNSIGNED NOT NULL,
    `corporation_id` int(10) UNSIGNED NOT NULL,
    `ship_type_id` int(10) UNSIGNED NOT NULL,
    `charisma` int(10) UNSIGNED NOT NULL,
    `intelligence` int(10) UNSIGNED NOT NULL,
    `memory` int(10) UNSIGNED NOT NULL,
    `perception` int(10) UNSIGNED NOT NULL,
    `willpower` int(10) UNSIGNED NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `bloodlines_race_id_idx` (`race_id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.characters
CREATE TABLE IF NOT EXISTS `characters` (
    `id` bigint(20) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `corporation_id` int(10) UNSIGNED NOT NULL,
    `alliance_id` int(10) UNSIGNED DEFAULT NULL,
    `faction_id` int(10) UNSIGNED DEFAULT NULL,
    `security_status` float DEFAULT NULL,
    `gender` enum('male', 'female') COLLATE utf8mb4_unicode_ci NOT NULL,
    `birthday` datetime NOT NULL,
    `title` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `bloodline_id` int(10) UNSIGNED NOT NULL,
    `race_id` int(10) UNSIGNED NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `characters_corporation_id_idx` (`corporation_id`) USING BTREE,
    KEY `characters_alliance_id_idx` (`alliance_id`) USING BTREE,
    KEY `characters_faction_id_idx` (`faction_id`) USING BTREE,
    KEY `characters_bloodline_id_idx` (`bloodline_id`) USING BTREE,
    KEY `characters_race_id_idx` (`race_id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.character_attributes
CREATE TABLE IF NOT EXISTS `character_attributes` (
    `character_id` bigint(20) UNSIGNED NOT NULL,
    `charisma` tinyint(3) UNSIGNED NOT NULL,
    `intelligence` tinyint(3) UNSIGNED NOT NULL,
    `memory` tinyint(3) UNSIGNED NOT NULL,
    `perception` tinyint(3) UNSIGNED NOT NULL,
    `willpower` tinyint(3) UNSIGNED NOT NULL,
    `bonus_remaps` tinyint(3) UNSIGNED NOT NULL DEFAULT '0',
    `last_remap_date` datetime DEFAULT NULL,
    `accrued_remap_cooldown_date` datetime DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`character_id`) USING BTREE,
    CONSTRAINT `character_skill_attributes_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `users` (`character_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.character_clone_meta
CREATE TABLE IF NOT EXISTS `character_clone_meta` (
    `character_id` bigint(20) UNSIGNED NOT NULL,
    `last_clone_jump_date` timestamp NULL DEFAULT NULL,
    `last_station_change_date` timestamp NULL DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`character_id`) USING BTREE,
    CONSTRAINT `character_clone_meta_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `users` (`character_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.character_corporation_history
CREATE TABLE IF NOT EXISTS `character_corporation_history` (
    `id` bigint(20) UNSIGNED NOT NULL,
    `record_id` int(10) UNSIGNED NOT NULL,
    `corporation_id` int(10) UNSIGNED NOT NULL,
    `is_deleted` tinyint(4) NOT NULL DEFAULT '0',
    `start_date` datetime NOT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`id`, `record_id`) USING BTREE,
    CONSTRAINT `character_corporation_history_id_foreign` FOREIGN KEY (`id`) REFERENCES `characters` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.character_home_clone
CREATE TABLE IF NOT EXISTS `character_home_clone` (
    `character_id` bigint(20) UNSIGNED NOT NULL,
    `location_id` bigint(20) UNSIGNED NOT NULL,
    `location_type` enum('station', 'structure') COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`character_id`) USING BTREE,
    CONSTRAINT `character_home_clone_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `users` (`character_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.character_implants
CREATE TABLE IF NOT EXISTS `character_implants` (
    `character_id` bigint(20) UNSIGNED NOT NULL,
    `implant_id` int(10) UNSIGNED NOT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`character_id`, `implant_id`) USING BTREE,
    KEY `character_implants_implant_id_foreign` (`implant_id`),
    CONSTRAINT `character_implants_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `users` (`character_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.character_jump_clones
CREATE TABLE IF NOT EXISTS `character_jump_clones` (
    `character_id` bigint(20) UNSIGNED NOT NULL,
    `jump_clone_id` int(10) UNSIGNED NOT NULL,
    `location_id` bigint(20) UNSIGNED NOT NULL,
    `location_type` enum('station', 'structure') COLLATE utf8mb4_unicode_ci NOT NULL,
    `implants` json NOT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`character_id`, `jump_clone_id`) USING BTREE,
    CONSTRAINT `character_jump_clones_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `users` (`character_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.character_skillqueue
CREATE TABLE IF NOT EXISTS `character_skillqueue` (
    `character_id` bigint(20) UNSIGNED NOT NULL,
    `queue_position` tinyint(3) UNSIGNED NOT NULL,
    `skill_id` int(10) UNSIGNED NOT NULL,
    `start_date` timestamp NULL DEFAULT NULL,
    `finish_date` timestamp NULL DEFAULT NULL,
    `finished_level` tinyint(3) UNSIGNED NOT NULL,
    `training_start_sp` int(10) UNSIGNED DEFAULT NULL,
    `level_start_sp` int(10) UNSIGNED DEFAULT NULL,
    `level_end_sp` int(10) UNSIGNED DEFAULT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`character_id`, `queue_position`) USING BTREE,
    KEY `character_skillqueue_skill_id_idx` (`skill_id`),
    KEY `character_skillqueue_start_date_idx` (`start_date`),
    KEY `character_skillqueue_finish_date_idx` (`finish_date`),
    CONSTRAINT `character_skillqueue_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `users` (`character_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.character_skills
CREATE TABLE IF NOT EXISTS `character_skills` (
    `character_id` bigint(20) UNSIGNED NOT NULL,
    `skill_id` int(10) UNSIGNED NOT NULL,
    `active_skill_level` tinyint(3) UNSIGNED NOT NULL,
    `skillpoints_in_skill` int(10) UNSIGNED NOT NULL,
    `trained_skill_level` tinyint(3) UNSIGNED NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`character_id`, `skill_id`) USING BTREE,
    CONSTRAINT `character_skills_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `users` (`character_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.character_skill_meta
CREATE TABLE IF NOT EXISTS `character_skill_meta` (
    `character_id` bigint(20) UNSIGNED NOT NULL,
    `total_sp` int(10) UNSIGNED NOT NULL,
    `unallocated_sp` int(10) UNSIGNED DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`character_id`) USING BTREE,
    CONSTRAINT `character_skill_meta_character_id_foreign` FOREIGN KEY (`character_id`) REFERENCES `users` (`character_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.corporations
CREATE TABLE IF NOT EXISTS `corporations` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `ticker` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `ceo_id` int(10) UNSIGNED NOT NULL,
    `creator_id` int(10) UNSIGNED NOT NULL,
    `alliance_id` int(10) UNSIGNED DEFAULT NULL,
    `home_station_id` int(10) UNSIGNED DEFAULT NULL,
    `faction_id` int(10) UNSIGNED DEFAULT NULL,
    `member_count` int(10) UNSIGNED NOT NULL,
    `shares` bigint(20) UNSIGNED DEFAULT NULL,
    `tax_rate` float NOT NULL,
    `url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `war_eligible` tinyint(4) NOT NULL DEFAULT '0',
    `date_founded` timestamp NULL DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `corporations_ceo_id_idx` (`ceo_id`),
    KEY `corporations_creator_id_idx` (`creator_id`),
    KEY `corporations_alliance_id_idx` (`alliance_id`),
    KEY `corporations_home_station_id_idx` (`home_station_id`),
    KEY `corporations_faction_id_idx` (`faction_id`),
    KEY `corporations_member_count_idx` (`member_count`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.etags
CREATE TABLE IF NOT EXISTS `etags` (
    `path` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL,
    `etag` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `cached_until` datetime NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`path`) USING BTREE,
    UNIQUE KEY `etags_path_unqiue_idx` (`path`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.factions
CREATE TABLE IF NOT EXISTS `factions` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `is_unique` tinyint(3) UNSIGNED NOT NULL DEFAULT '0',
    `size_factor` float NOT NULL,
    `station_count` int(10) UNSIGNED NOT NULL,
    `station_system_count` int(10) UNSIGNED NOT NULL,
    `corporation_id` int(10) UNSIGNED DEFAULT NULL,
    `militia_corporation_id` int(10) UNSIGNED DEFAULT NULL,
    `solar_system_id` int(10) UNSIGNED DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `factions_solar_system_id_idx` (`solar_system_id`),
    KEY `factions_corporation_id_idx` (`corporation_id`),
    KEY `factions_militia_corporation_id_idx` (`militia_corporation_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.map_constellations
CREATE TABLE IF NOT EXISTS `map_constellations` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `region_id` int(10) UNSIGNED NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `constellations_region_id_idx` (`region_id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.map_regions
CREATE TABLE IF NOT EXISTS `map_regions` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.map_solar_systems
CREATE TABLE IF NOT EXISTS `map_solar_systems` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `constellation_id` int(10) UNSIGNED NOT NULL,
    `security_class` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `security_status` float NOT NULL,
    `star_id` int(10) UNSIGNED DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `solar_systems_constellation_id_idx` (`constellation_id`) USING BTREE,
    KEY `solar_systems_star_id_idx` (`star_id`) USING BTREE,
    KEY `solar_systems_security_status` (`security_status`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.map_stations
CREATE TABLE IF NOT EXISTS `map_stations` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `system_id` int(10) UNSIGNED NOT NULL,
    `type_id` int(10) UNSIGNED NOT NULL,
    `race_id` int(10) UNSIGNED DEFAULT NULL,
    `owner_corporation_id` int(10) UNSIGNED DEFAULT NULL,
    `max_dockable_ship_volume` float NOT NULL,
    `office_rental_cost` float NOT NULL,
    `reprocessing_efficiency` float NOT NULL,
    `reprocessing_stations_take` float NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `stations_system_id_idx` (`system_id`) USING BTREE,
    KEY `stations_type_id_idx` (`type_id`) USING BTREE,
    KEY `stations_race_id_idx` (`race_id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.map_structures
CREATE TABLE IF NOT EXISTS `map_structures` (
    `id` bigint(20) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `solar_system_id` int(10) UNSIGNED NOT NULL,
    `type_id` int(10) UNSIGNED NOT NULL,
    `owner_id` int(10) UNSIGNED NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `structures_solar_system_id_idx` (`solar_system_id`) USING BTREE,
    KEY `structures_type_id_idx` (`type_id`) USING BTREE,
    KEY `structures_owner_id_idx` (`owner_id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.races
CREATE TABLE IF NOT EXISTS `races` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.types
CREATE TABLE IF NOT EXISTS `types` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `group_id` int(10) UNSIGNED NOT NULL,
    `published` tinyint(3) UNSIGNED NOT NULL DEFAULT '0',
    `capacity` float NOT NULL DEFAULT '0',
    `market_group_id` int(10) UNSIGNED DEFAULT NULL,
    `mass` float NOT NULL DEFAULT '0',
    `packaged_volume` float NOT NULL DEFAULT '0',
    `portion_size` int(11) DEFAULT NULL,
    `radius` float DEFAULT NULL,
    `volume` float NOT NULL DEFAULT '0',
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `types_group_id_idx` (`group_id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.type_attributes
CREATE TABLE IF NOT EXISTS `type_attributes` (
    `type_id` int(10) UNSIGNED NOT NULL,
    `attribute_id` int(10) UNSIGNED NOT NULL,
    `value` float NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`type_id`, `attribute_id`) USING BTREE,
    CONSTRAINT `types_attributes_type_id_types_id_foreign` FOREIGN KEY (`type_id`) REFERENCES `types` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.type_categories
CREATE TABLE IF NOT EXISTS `type_categories` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `published` tinyint(4) NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.type_groups
CREATE TABLE IF NOT EXISTS `type_groups` (
    `id` int(10) UNSIGNED NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `published` tinyint(4) NOT NULL,
    `category_id` int(10) UNSIGNED NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `type_groups_category_id_idx` (`category_id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Data exporting was unselected.
-- Dumping structure for table skillz.users
CREATE TABLE IF NOT EXISTS `users` (
    `id` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
    `character_id` bigint(20) UNSIGNED NOT NULL,
    `access_token` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `refresh_token` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `expires` datetime NOT NULL,
    `owner_hash` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
    `scopes` json NOT NULL,
    `disabled` tinyint(3) UNSIGNED NOT NULL DEFAULT '0',
    `disabled_reason` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `disabled_timestamp` datetime DEFAULT NULL,
    `last_login` datetime NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `character_id` (`character_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;