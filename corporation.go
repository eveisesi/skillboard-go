package skillz

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type CorporationRepository interface {
	Corporation(ctx context.Context, corporationID uint) (*Corporation, error)
	CreateCorporation(ctx context.Context, corporation *Corporation) error
	UpdateCorporation(ctx context.Context, corporation *Corporation) error
	CorporationAllianceHistory(ctx context.Context, corporationID uint) ([]*CorporationAllianceHistory, error)
	CreateCorporationAllianceHistory(ctx context.Context, records []*CorporationAllianceHistory) ([]*CorporationAllianceHistory, error)
}

type Corporation struct {
	ID            uint        `db:"id,omitempty" json:"id"`
	AllianceID    null.Uint   `db:"alliance_id,omitempty" json:"alliance_id,omitempty"`
	CeoID         uint        `db:"ceo_id" json:"ceo_id"`
	CreatorID     uint        `db:"creator_id" json:"creator_id"`
	DateFounded   null.Time   `db:"date_founded,omitempty" json:"date_founded,omitempty"`
	FactionID     null.Uint   `db:"faction_id,omitempty" json:"faction_id,omitempty"`
	HomeStationID null.Uint   `db:"home_station_id,omitempty" json:"home_station_id,omitempty"`
	MemberCount   uint        `db:"member_count" json:"member_count"`
	Name          string      `db:"name" json:"name"`
	Shares        uint64      `db:"shares,omitempty" json:"shares,omitempty"`
	TaxRate       float32     `db:"tax_rate" json:"tax_rate"`
	Ticker        string      `db:"ticker" json:"ticker"`
	URL           null.String `db:"url,omitempty" json:"url,omitempty"`
	WarEligible   bool        `db:"war_eligible" json:"war_eligible"`
	CreatedAt     time.Time   `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt     time.Time   `db:"updated_at" json:"updated_at" deep:"-"`
}

type CorporationAllianceHistory struct {
	CorporationID uint      `db:"corporation_id" json:"id"`
	RecordID      uint      `db:"record_id" json:"record_id"`
	AllianceID    null.Uint `db:"alliance_id" json:"alliance_id"`
	IsDeleteed    null.Bool `db:"is_deleted" json:"is_deleted"`
	StartDate     time.Time `db:"start_date" json:"start_date"`
	CreatedAt     time.Time `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at" deep:"-"`
}
