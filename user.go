package skillz

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
)

type User struct {
	ID                uuid.UUID    `db:"id" json:"id"`
	CharacterID       uint         `db:"main_id,omitempty" json:"main_id"`
	AccessToken       string       `db:"access_token" json:"access_token"`
	RefreshToken      string       `db:"refresh_token" json:"refresh_token"`
	Expires           time.Time    `db:"expires," json:"expires"`
	OwnerHash         string       `db:"owner_hash" json:"owner_hash"`
	Scopes            MemberScopes `db:"scopes,omitempty" json:"scopes,omitempty"`
	Disabled          bool         `db:"disabled" json:"disabled"`
	DisabledReason    null.String  `db:"disabled_reason,omitempty" json:"disabled_reason"`
	DisabledTimestamp null.Time    `db:"disabled_timestamp,omitempty" json:"disabled_timestamp"`
	LastLogin         time.Time    `db:"last_login" json:"last_login"`
	CreatedAt         time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time    `db:"updated_at" json:"updated_at"`
}
