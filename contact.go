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
	CharacterID uint64      `db:"character_id"`
	ContactID   uint        `db:"contact_id"`
	ContactType ContactType `db:"contact_type"`
	Standing    float64     `db:"standing"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`
}

type ContactType string

const (
	CharacterContactType   ContactType = "character"
	CorporationContactType ContactType = "corporation"
	AllianceContactType    ContactType = "alliance"
	FactionContactType     ContactType = "faction"
)
