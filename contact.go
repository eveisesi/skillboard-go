package skillz

import (
	"context"
	"time"
)

type ContactRepository interface {
	CharacterContacts(ctx context.Context, characterID uint64) ([]*CharacterContact, error)
	CreateCharacterContacts(ctx context.Context, contacts []*CharacterContact) error
	DeleteCharacterContacts(ctx context.Context, characterID uint64) error
}

type CharacterContact struct {
	CharacterID uint64      `db:"character_id" json:"character_id"`
	ContactID   uint        `db:"contact_id" json:"contact_id"`
	ContactType ContactType `db:"contact_type" json:"contact_type"`
	Standing    float64     `db:"standing" json:"standing"`
	CreatedAt   time.Time   `db:"created_at" json:"-"`
	UpdatedAt   time.Time   `db:"updated_at" json:"-"`

	Alliance    *Alliance    `json:"alliance,omitempty"`
	Corporation *Corporation `json:"corporation,omitempty"`
	Character   *Character   `json:"character,omitempty"`
}

type ContactType string

const (
	CharacterContactType   ContactType = "character"
	CorporationContactType ContactType = "corporation"
	AllianceContactType    ContactType = "alliance"
	FactionContactType     ContactType = "faction"
)
