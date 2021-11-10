package skillz

import (
	"context"
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
)

type ScopeMap map[Scope][]ScopeResolver

type ScopeResolver struct {
	Name string
	Func func(context.Context, *User) (*Etag, error)
}

type Scope string

const (
	ReadImplantsV1   Scope = "esi-clones.read_implants.v1"
	ReadLocationV1   Scope = "esi-location.read_location.v1"
	ReadOnlineV1     Scope = "esi-location.read_online.v1"
	ReadShipV1       Scope = "esi-location.read_ship_type.v1"
	ReadSkillQueueV1 Scope = "esi-skills.read_skillqueue.v1"
	ReadSkillsV1     Scope = "esi-skills.read_skills.v1"
)

var AllScopes = []Scope{
	ReadImplantsV1, ReadLocationV1,
	ReadOnlineV1, ReadShipV1,
	ReadSkillQueueV1, ReadSkillsV1,
}

func (s Scope) String() string {
	return string(s)
}

type MemberScope struct {
	Scope  Scope     `db:"scope" json:"scope"`
	Expiry null.Time `db:"expiry,omitempty" json:"expiry,omitempty"`
}

type MemberScopes []MemberScope

func (s *MemberScopes) Scan(value interface{}) error {

	switch data := value.(type) {
	case []byte:
		var scopes MemberScopes
		err := json.Unmarshal(data, &scopes)
		if err != nil {
			return err
		}

		*s = scopes
	}

	return nil
}

func (s MemberScopes) Value() (driver.Value, error) {

	if len(s) == 0 {
		return `[]`, nil
	}
	data, err := json.Marshal(s)

	return data, errors.Wrap(err, "[MemberScopes] Failed to marshal scope for storage in data store")
}
