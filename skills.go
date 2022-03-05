package skillz

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type CharacterSkillRepository interface {
	memberAttributesRepository
	memberSkillsRepository
	memberSkillQueueRepository
	memberFlyableShipRepository
}

type memberAttributesRepository interface {
	CharacterAttributes(ctx context.Context, characterID uint64) (*CharacterAttributes, error)
	CreateCharacterAttributes(ctx context.Context, attributes *CharacterAttributes) error
	DeleteCharacterAttributes(ctx context.Context, characterID uint64) error
}

type memberSkillsRepository interface {
	CharacterSkillMeta(ctx context.Context, characterID uint64) (*CharacterSkillMeta, error)
	CreateCharacterSkillMeta(ctx context.Context, meta *CharacterSkillMeta) error
	DeleteCharacterSkillMeta(ctx context.Context, characterID uint64) error

	CharacterSkills(ctx context.Context, characterID uint64) ([]*CharacterSkill, error)
	CreateCharacterSkills(ctx context.Context, skills []*CharacterSkill) error
	DeleteCharacterSkills(ctx context.Context, characterID uint64) error
}

type memberSkillQueueRepository interface {
	CharacterSkillQueue(ctx context.Context, characterID uint64) ([]*CharacterSkillQueue, error)
	CreateCharacterSkillQueue(ctx context.Context, positions []*CharacterSkillQueue) error
	DeleteCharacterSkillQueue(ctx context.Context, characterID uint64) error
}

type memberFlyableShipRepository interface {
	CharacterFlyableShips(ctx context.Context, characterID uint64) ([]*CharacterFlyableShip, error)
	CreateCharacterFlyableShips(ctx context.Context, ships []*CharacterFlyableShip) error
	DeleteCharacterFlyableShips(ctx context.Context, characterID uint64) error
}

type CharacterAttributes struct {
	CharacterID              uint64    `db:"character_id" json:"character_id"`
	Charisma                 uint      `db:"charisma" json:"charisma"`
	Intelligence             uint      `db:"intelligence" json:"intelligence"`
	Memory                   uint      `db:"memory" json:"memory"`
	Perception               uint      `db:"perception" json:"perception"`
	Willpower                uint      `db:"willpower" json:"willpower"`
	BonusRemaps              null.Uint `db:"bonus_remaps,omitempty" json:"bonus_remaps,omitempty"`
	LastRemapDate            null.Time `db:"last_remap_date,omitempty" json:"last_remap_date,omitempty"`
	AccruedRemapCooldownDate null.Time `db:"accrued_remap_cooldown_date,omitempty" json:"accrued_remap_cooldown_date,omitempty"`
	CreatedAt                time.Time `db:"created_at" json:"-" deep:"-"`
	UpdatedAt                time.Time `db:"updated_at" json:"-" deep:"-"`
}

type CharacterSkillQueueSummary struct {
	Summary []*QueueGroupSummary   `json:"summary"`
	Queue   []*CharacterSkillQueue `json:"queue"`
}

type QueueGroupSummary struct {
	Group       *Group        `json:"group"`
	Count       uint          `json:"count"`
	Skillpoints uint          `json:"skillpoints"`
	Duration    time.Duration `json:"duration"`
}

type CharacterSkillQueue struct {
	CharacterID     uint64    `db:"character_id" json:"character_id" deep:"-"`
	QueuePosition   uint      `db:"queue_position" json:"queue_position"`
	SkillID         uint      `db:"skill_id" json:"skill_id"`
	FinishedLevel   uint      `db:"finished_level" json:"finished_level"`
	TrainingStartSp null.Uint `db:"training_start_sp,omitempty" json:"training_start_sp,omitempty"`
	LevelStartSp    null.Uint `db:"level_start_sp,omitempty" json:"level_start_sp,omitempty"`
	LevelEndSp      null.Uint `db:"level_end_sp,omitempty" json:"level_end_sp,omitempty"`
	StartDate       null.Time `db:"start_date,omitempty" json:"start_date,omitempty"`
	FinishDate      null.Time `db:"finish_date,omitempty" json:"finish_date,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"-" deep:"-"`

	Type *Type `json:"info"`
}

type CharacterSkillMeta struct {
	CharacterID   uint64    `db:"character_id" json:"character_id"`
	TotalSP       uint      `db:"total_sp" json:"total_sp"`
	UnallocatedSP null.Uint `db:"unallocated_sp,omitempty" json:"unallocated_sp,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	UpdatedAt     time.Time `db:"updated_at" json:"-"`

	Skills []*CharacterSkill `json:"skills,omitempty"`
}

type CharacterSkill struct {
	CharacterID        uint64    `db:"character_id" json:"character_id"`
	ActiveSkillLevel   uint      `db:"active_skill_level" json:"active_skill_level"`
	SkillID            uint      `db:"skill_id" json:"skill_id"`
	SkillpointsInSkill uint      `db:"skillpoints_in_skill" json:"skillpoints_in_skill"`
	TrainedSkillLevel  uint      `db:"trained_skill_level" json:"trained_skill_level"`
	CreatedAt          time.Time `db:"created_at" json:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"-"`

	Info *Type `json:"type,omitempty"`
}

type CharacterSkillGroup struct {
	*SkillGroup
	TotalGroupSP uint `json:"totalGroupSP"`
}

type SkillGroup struct {
	*Group
	Skills []*SkillType `json:"skills"`
}

type SkillType struct {
	*Type
	Skill *CharacterSkill     `json:"skill"`
	Rank  *TypeDogmaAttribute `json:"rank"`
}

// CharacterFlyableShip is a model of the database table
type CharacterFlyableShip struct {
	CharacterID uint64    `db:"character_id" json:"character_id"`
	ShipTypeID  uint      `db:"ship_type_id" json:"ship_type_id"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
}

// ShipGroup is what is used on the frontend to render the Flyable Ships view
type ShipGroup struct {
	*Group
	Ships []*ShipType
}

type ShipType struct {
	*Type
	Flyable bool
}
