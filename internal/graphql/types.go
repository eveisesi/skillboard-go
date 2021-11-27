package graphql

import (
	"time"

	"github.com/eveisesi/skillz"
	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
)

type Skillboard struct {
	User *User `json:"user"`
}

type User struct {
	ID          uuid.UUID `json:"id"`
	CharacterID uint64    `json:"characterID"`
	Scopes      []string  `json:"scopes"`
	IsNew       bool      `json:"isNew"`
	LastLogin   time.Time `json:"lastLogin"`

	Character *Character    `json:"character"`
	SkillMeta *SkillMeta    `json:"skillMeta"`
	Skills    []*SkillGroup `json:"skills"`
	// Skills     *Skills     `json:"skills"`
	Implants   []*Implant  `json:"implants"`
	Attributes *Attributes `json:"attributes"`
	Queue      []*Queue    `json:"queue"`
}

type Alliance struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Corporation struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`

	Alliance *Alliance `json:"alliance"`
}

type Character struct {
	ID             uint64       `db:"id" json:"id"`
	Name           string       `db:"name" json:"name"`
	SecurityStatus null.Float64 `db:"securityStatus" json:"securityStatus"`
	Birthday       time.Time    `db:"birthday" json:"birthday"`

	Corporation *Corporation `json:"corporation"`
}

type SkillMeta struct {
	CharacterID   uint64    `json:"characterID"`
	TotalSP       uint      `json:"totalSp"`
	UnallocatedSP null.Uint `json:"unallocatedSp"`
}

type SkillGroup struct {
	TotalGroupSP uint `json:"totalGroupSP"`
	Info         *InvGroup
	Skills       []*Skill
}

type Skill struct {
	CharacterID        uint64 `json:"characterID"`
	ActiveSkillLevel   uint   `json:"activeSkillLevel"`
	SkillID            uint   `json:"skillID"`
	SkillpointsInSkill uint   `json:"skillpointsInSkill"`
	TrainedSkillLevel  uint   `json:"trainedSkillLevel"`

	Info *InvType `json:"info"`
}

type Implant struct {
	CharacterID uint64 `json:"characterID"`
	ImplantID   uint   `json:"implantID"`

	Implant *skillz.Type `json:"implant"`
}

type Attributes struct {
	CharacterID              uint64    `json:"characterID"`
	Charisma                 uint      `json:"charisma"`
	Intelligence             uint      `json:"intelligence"`
	Memory                   uint      `json:"memory"`
	Perception               uint      `json:"perception"`
	Willpower                uint      `json:"willpower"`
	BonusRemaps              null.Uint `json:"bonusRemaps"`
	LastRemapDate            null.Time `json:"lastRemapDate"`
	AccruedRemapCooldownDate null.Time `json:"accruedRemapCooldownDate"`
}

type Queue struct {
	CharacterID     uint64    `json:"characterID"`
	QueuePosition   uint      `json:"queuePosition"`
	SkillID         uint      `json:"skillID"`
	FinishedLevel   uint      `json:"finishedLevel"`
	TrainingStartSp null.Uint `json:"trainingStartSp,omitempty"`
	LevelStartSp    null.Uint `json:"levelStartSp,omitempty"`
	LevelEndSp      null.Uint `json:"levelEndSp,omitempty"`
	StartDate       null.Time `json:"startDate,omitempty"`
	FinishDate      null.Time `json:"finishDate,omitempty"`

	Info *skillz.Type
}

type InvType struct {
	ID   uint   `db:"id" json:"id"`
	Name string `db:"name" json:"name"`

	Attribute []*InvTypeAttribute `json:"attributes"`
	Group     *InvGroup           `json:"group"`
}

type InvTypeAttribute struct {
	AttributeID uint    `json:"attributeID"`
	Value       float64 `json:"value"`
}

type InvGroup struct {
	ID   uint   `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
