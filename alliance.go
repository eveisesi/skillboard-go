package skillz

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type AllianceRepository interface {
	Alliance(ctx context.Context, id uint) (*Alliance, error)
	CreateAlliance(ctx context.Context, alliance *Alliance) error
	UpdateAlliance(ctx context.Context, alliance *Alliance) error
}

// Alliance is an object representing the database table.
type Alliance struct {
	ID                    uint      `db:"id" json:"id"`
	Name                  string    `db:"name" json:"name"`
	Ticker                string    `db:"ticker" json:"ticker"`
	DateFounded           time.Time `db:"date_founded" json:"date_founded"`
	CreatorID             uint      `db:"creator_id" json:"creator_id"`
	CreatorCorporationID  uint      `db:"creator_corporation_id" json:"creator_corporation_id"`
	ExecutorCorporationID uint      `db:"executor_corporation_id" json:"executor_corporation_id"`
	FactionID             null.Uint `db:"faction_id,omitempty" json:"faction_id,omitempty"`
	IsClosed              bool      `db:"is_closed" json:"is_closed"`
	CreatedAt             time.Time `db:"created_at" json:"-"`
	UpdatedAt             time.Time `db:"updated_at" json:"-"`
}
