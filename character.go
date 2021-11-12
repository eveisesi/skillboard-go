package skillz

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type CharacterRepository interface {
	characterRepository
	characterHistoryRepository
}

type characterRepository interface {
	Character(ctx context.Context, id uint64) (*Character, error)
	// Characters(ctx context.Context, operators ...*Operator) ([]*Character, error)
	CreateCharacter(ctx context.Context, character *Character) error
	UpdateCharacter(ctx context.Context, character *Character) error
}

type characterHistoryRepository interface {
	CharacterCorporationHistory(ctx context.Context, characterID uint64) ([]*CharacterCorporationHistory, error)
	CreateCharacterCorporationHistory(ctx context.Context, records []*CharacterCorporationHistory) ([]*CharacterCorporationHistory, error)
}

type Character struct {
	ID             uint64       `db:"id" json:"id"`
	Name           string       `db:"name" json:"name"`
	CorporationID  uint         `db:"corporation_id" json:"corporation_id"`
	AllianceID     null.Uint    `db:"alliance_id,omitempty" json:"alliance_id,omitempty"`
	FactionID      null.Uint    `db:"faction_id,omitempty" json:"faction_id,omitempty"`
	SecurityStatus null.Float64 `db:"security_status,omitempty" json:"security_status,omitempty"`
	Gender         string       `db:"gender" json:"gender"`
	Birthday       time.Time    `db:"birthday" json:"birthday"`
	Title          null.String  `db:"title,omitempty" json:"title,omitempty"`
	BloodlineID    uint         `db:"bloodline_id" json:"bloodline_id"`
	RaceID         uint         `db:"race_id" json:"race_id"`
	CreatedAt      time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time    `db:"updated_at" json:"updated_at"`
}

type CharacterCorporationHistory struct {
	CharacterID   uint64    `db:"character_id" json:"character_id"`
	RecordID      uint64    `db:"record_id" json:"record_id"`
	CorporationID uint      `db:"corporation_id" json:"corporation_id"`
	IsDeleted     bool      `db:"is_deleted" json:"is_deleted"`
	StartDate     time.Time `db:"start_date" json:"start_date"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
