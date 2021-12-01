package mysql

const (
	ColumnCharacterID string = "character_id"
	ColumnCreatedAt   string = "created_at"
	ColumnUpdatedAt   string = "updated_at"
)

const (
	TableAlliances                   string = "alliances"
	TableBloodlines                  string = "bloodlines"
	TableCharacters                  string = "characters"
	TableCharacterAttributes         string = "character_attributes"
	TableCharacterCloneMeta          string = "character_clone_meta"
	TableCharacterCorporationHistory string = "character_corporation_history"
	TableCharacterHomeClone          string = "character_home_clone"
	TableCharacterImplants           string = "character_implants"
	TableCharacterJumpClones         string = "character_jump_clones"
	TableCharacterSkillQueue         string = "character_skillqueue"
	TableCharacterSkills             string = "character_skills"
	TableCharacterSkillMeta          string = "character_skill_meta"
	TableCorporations                string = "corporations"
	TableCorporationAllianceHistory  string = "corporation_alliance_history"
	TableEtags                       string = "etags"
	TableFactions                    string = "factions"
	TableMapConstellations           string = "map_constellations"
	TableMapRegions                  string = "map_regions"
	TableMapSolarSystems             string = "map_solar_systems"
	TableMapStations                 string = "map_stations"
	TableStructures                  string = "map_structures"
	TableRaces                       string = "races"
	TableTypes                       string = "types"
	TableTypeAttributes              string = "type_attributes"
	TableTypeCategories              string = "type_categories"
	TableTypeGroups                  string = "type_groups"
	TableUsers                       string = "users"
)

const (
	prefixFormat string = "[%s.%s]"
	errorFFormat string = "[%s.%s] %s"
)
const (
	allianceRepositoryIdentifier    string = "AllianceRepository"
	characterRepositoryIdentifier   string = "CharacterRepository"
	cloneRepositoryIdentifier       string = "CloneRepository"
	corporationRepositoryIdentifier string = "CorporationRepository"
	etagRepositoryIdentifier        string = "ETagRepository"
	skillsRepositoryIdentifier      string = "SkillsRepository"
	universeRepositoryIdentifier    string = "UniverseRepository"
	userRepositoryIdentifier        string = "UserRepository"
)
