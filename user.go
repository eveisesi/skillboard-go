package skillz

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
	"golang.org/x/oauth2"
)

type UserRepository interface {
	User(ctx context.Context, id uuid.UUID) (*User, error)
	UserByCharacterID(ctx context.Context, characterID uint64) (*User, error)
	SearchUsers(ctx context.Context, q string) ([]*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error

	UserSettings(ctx context.Context, id uuid.UUID) (*UserSettings, error)
	CreateUserSettings(ctx context.Context, settings *UserSettings) error

	UsersSortedByProcessedAtLimit(ctx context.Context, limit uint64) ([]*User, error)

	NewUsersBySP(ctx context.Context) ([]*User, error)
}

type User struct {
	ID                uuid.UUID   `db:"id" json:"id"`
	CharacterID       uint64      `db:"character_id,omitempty" json:"character_id"`
	AccessToken       string      `db:"access_token" json:"-"`
	RefreshToken      string      `db:"refresh_token" json:"-"`
	Expires           time.Time   `db:"expires," json:"expires"`
	OwnerHash         string      `db:"owner_hash" json:"owner_hash"`
	Scopes            UserScopes  `db:"scopes,omitempty" json:"scopes,omitempty"`
	IsNew             bool        `db:"is_new" json:"is_new"`
	Disabled          bool        `db:"disabled" json:"disabled"`
	DisabledReason    null.String `db:"disabled_reason,omitempty" json:"disabled_reason"`
	DisabledTimestamp null.Time   `db:"disabled_timestamp,omitempty" json:"disabled_timestamp"`
	LastLogin         time.Time   `db:"last_login" json:"last_login"`
	LastProcessed     null.Time   `db:"last_processed" json:"last_processed"`
	CreatedAt         time.Time   `db:"created_at" json:"-"`
	UpdatedAt         time.Time   `db:"updated_at" json:"-"`

	Errors []error `json:"errors,omitempty"`

	Character  *Character              `json:"character,omitempty"`
	Settings   *UserSettings           `json:"settings,omitempty"`
	Skills     []*CharacterSkill       `json:"skillz,omitempty"`
	Queue      []*CharacterSkillQueue  `json:"queue,omitempty"`
	Attributes *CharacterAttributes    `json:"attributes,omitempty"`
	Flyable    []*CharacterFlyableShip `json:"flyable,omitempty"`
}

func (i *User) ApplyToken(t *oauth2.Token) {
	i.AccessToken = t.AccessToken
	i.RefreshToken = t.RefreshToken
	i.Expires = t.Expiry
}

type UserSettings struct {
	UserID        uuid.UUID `db:"user_id" json:"-"`
	HideClones    bool      `db:"hide_clones" json:"hide_clones"`
	HideQueue     bool      `db:"hide_queue" json:"hide_queue"`
	HideStandings bool      `db:"hide_standings" json:"hide_standings"`
	HideShips     bool      `db:"hide_ships" json:"hide_ships"`
	ForSale       bool      `db:"for_sale" json:"for_sale"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	UpdatedAt     time.Time `db:"updated_at" json:"-"`
}

type UserSearchResult struct {
	*User
	Info *Character `json:"info"`
}

type UserWithSkillMeta struct {
	*User
	Meta   *CharacterSkillMeta    `json:"meta"`
	Skills []*CharacterSkill      `json:"skills"`
	Queue  []*CharacterSkillQueue `json:"skillQueue"`
	Info   *Character             `json:"info"`
}
