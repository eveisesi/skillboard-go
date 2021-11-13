package skillz

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
)

type UserRepository interface {
	User(ctx context.Context, id uuid.UUID) (*User, error)
	UserByCharacterID(ctx context.Context, characterID uint64) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
}

type User struct {
	ID                uuid.UUID   `db:"id" json:"id"`
	CharacterID       uint64      `db:"character_id,omitempty" json:"character_id"`
	AccessToken       string      `db:"access_token" json:"access_token"`
	RefreshToken      string      `db:"refresh_token" json:"refresh_token"`
	Expires           time.Time   `db:"expires," json:"expires"`
	OwnerHash         string      `db:"owner_hash" json:"owner_hash"`
	Scopes            UserScopes  `db:"scopes,omitempty" json:"scopes,omitempty"`
	Disabled          bool        `db:"disabled" json:"disabled"`
	DisabledReason    null.String `db:"disabled_reason,omitempty" json:"disabled_reason"`
	DisabledTimestamp null.Time   `db:"disabled_timestamp,omitempty" json:"disabled_timestamp"`
	LastLogin         time.Time   `db:"last_login" json:"last_login"`
	CreatedAt         time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time   `db:"updated_at" json:"updated_at"`
}
