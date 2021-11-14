package skillz

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/volatiletech/null"
)

type CloneRepository interface {
	CharacterCloneMeta(ctx context.Context, characterID uint64) (*CharacterCloneMeta, error)
	CreateCharacterCloneMeta(ctx context.Context, meta *CharacterCloneMeta) error
	UpdateCharacterCloneMeta(ctx context.Context, meta *CharacterCloneMeta) error

	CharacterDeathClone(ctx context.Context, characterID uint64) (*CharacterDeathClone, error)
	CreateCharacterDeathClone(ctx context.Context, death *CharacterDeathClone) error
	UpdateCharacterDeathClone(ctx context.Context, death *CharacterDeathClone) error

	CharacterJumpClones(ctx context.Context, characterID uint64) ([]*CharacterJumpClone, error)
	CreateCharacterJumpClones(ctx context.Context, clones []*CharacterJumpClone) error
	DeleteCharacterJumpClones(ctx context.Context, characterID uint64) error

	CharacterImplants(ctx context.Context, characterID uint64) ([]*CharacterImplant, error)
	CreateCharacterImplants(ctx context.Context, implants []*CharacterImplant) error
	DeleteCharacterImplants(ctx context.Context, characterID uint64) error
}

type CharacterCloneMeta struct {
	CharacterID           uint64    `db:"character_id" json:"character_id"`
	LastCloneJumpDate     null.Time `db:"last_clone_jump_date" json:"last_clone_jump_date"`
	LastStationChangeDate null.Time `db:"last_station_change_date" json:"last_station_change_date"`
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time `db:"updated_at" json:"updated_at"`

	HomeLocation *CharacterDeathClone  `json:"home_location"`
	JumpClones   []*CharacterJumpClone `json:"jump_clones"`
}

type CloneLocationType string

const (
	CloneLocationTypeStation   CloneLocationType = "station"
	CloneLocationTypeStructure CloneLocationType = "structure"
)

var AllCloneLocationTypes = []CloneLocationType{CloneLocationTypeStation, CloneLocationTypeStructure}

func (c CloneLocationType) Valid() bool {
	for _, t := range AllCloneLocationTypes {
		if c == t {
			return true
		}
	}

	return false
}

func (c CloneLocationType) String() string {
	return string(c)
}

type CharacterDeathClone struct {
	CharacterID  uint64            `db:"character_id" json:"character_id"`
	LocationID   uint64            `db:"location_id" json:"location_id"`
	LocationType CloneLocationType `db:"location_type" json:"location_type"`
	CreatedAt    time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time         `db:"updated_at" json:"updated_at"`
}

type CharacterJumpClone struct {
	CharacterID  uint64            `db:"character_id" json:"character_id"`
	JumpCloneID  uint              `db:"jump_clone_id" json:"jump_clone_id"`
	LocationID   uint64            `db:"location_id" json:"location_id"`
	LocationType CloneLocationType `db:"location_type" json:"location_type"`
	Implants     SliceUint         `db:"implants" json:"implants"`
	CreatedAt    time.Time         `db:"created_at" json:"created_at"`
}

type CharacterImplant struct {
	CharacterID uint64    `db:"character_id" json:"character_id"`
	ImplantID   uint      `db:"implant_id" json:"implant_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type SliceUint []uint64

func (s *SliceUint) Scan(value interface{}) error {

	switch data := value.(type) {
	case []byte:
		err := json.Unmarshal(data, s)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s SliceUint) Value() (driver.Value, error) {

	var data []byte
	var err error
	if len(s) == 0 {
		data, err = json.Marshal([]interface{}{})
	} else {
		data, err = json.Marshal(s)
	}
	if err != nil {
		return nil, fmt.Errorf("[SliceUint] Failed to marshal slice of uints for storage in data store: %w", err)
	}

	return data, nil

}

func (s SliceUint) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint64(s))
}

func (s *SliceUint) UnmarshalJSON(value []byte) error {

	x := make([]uint64, 0)
	err := json.Unmarshal(value, &x)
	if err != nil {
		return err
	}
	a := SliceUint(x)
	*s = a
	return nil

}
