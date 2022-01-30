package skillz

import (
	"context"
	"time"
)

type CloneRepository interface {
	CharacterImplants(ctx context.Context, characterID uint64) ([]*CharacterImplant, error)
	CreateCharacterImplants(ctx context.Context, implants []*CharacterImplant) error
	DeleteCharacterImplants(ctx context.Context, characterID uint64) error
}

type CharacterImplant struct {
	CharacterID uint64    `db:"character_id" json:"character_id"`
	ImplantID   uint      `db:"implant_id" json:"implant_id"`
	Slot        uint      `db:"slot"`
	CreatedAt   time.Time `db:"created_at" json:"-"`

	Type *Type `json:"info,omitempty"`
}
