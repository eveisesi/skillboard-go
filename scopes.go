package skillz

import (
	"context"
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

// How to use Processor interface
// Call Scopes to get a slice of scopes that the Processor supports
// Loops over that slice, if the user use one of more of the scopes
// that the processor supports, call the Process func passing in the
// User. That process func will evaluate the user's scopes to determine
// which internal functionality to call

type Processor interface {
	Process(ctx context.Context, user *User) error
	Scopes() []Scope
}

type ScopeProcessors []Processor

type ScopeResolver struct {
	Name string
	Func func(context.Context, *User) (*Etag, error)
}

type Scope string

const (
	ReadContactsV1   Scope = "esi-characters.read_contacts.v1"
	ReadClonesV1     Scope = "esi-clones.read_clones.v1"
	ReadImplantsV1   Scope = "esi-clones.read_implants.v1"
	ReadLocationV1   Scope = "esi-location.read_location.v1"
	ReadOnlineV1     Scope = "esi-location.read_online.v1"
	ReadShipV1       Scope = "esi-location.read_ship_type.v1"
	ReadSkillQueueV1 Scope = "esi-skills.read_skillqueue.v1"
	ReadSkillsV1     Scope = "esi-skills.read_skills.v1"
	ReadStructuresV1 Scope = "esi-universe.read_structures.v1"
)

var AllScopes = []Scope{
	ReadImplantsV1, ReadClonesV1, ReadContactsV1,
	ReadLocationV1, ReadOnlineV1, ReadShipV1,
	ReadSkillQueueV1, ReadSkillsV1, ReadStructuresV1,
}

func (s Scope) String() string {
	return string(s)
}

type UserScopes []Scope

func (s *UserScopes) Scan(value interface{}) error {

	switch data := value.(type) {
	case []byte:
		var scopes UserScopes
		err := json.Unmarshal(data, &scopes)
		if err != nil {
			return err
		}

		*s = scopes
	}

	return nil
}

func (s UserScopes) Value() (driver.Value, error) {

	if len(s) == 0 {
		return `[]`, nil
	}
	data, err := json.Marshal(s)

	return data, errors.Wrap(err, "[UserScopes] Failed to marshal scope for data store")

}
