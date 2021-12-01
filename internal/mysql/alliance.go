package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type allianceRepository struct {
	db       QueryExecContext
	alliance tableConf
}

const (
	AllianceID                    string = "id"
	AllianceName                  string = "name"
	AllianceTicker                string = "ticker"
	AllianceDateFounded           string = "date_founded"
	AllianceCreatorID             string = "creator_id"
	AllianceCreatorCorporationID  string = "creator_corporation_id"
	AllianceExecutorCorporationID string = "executor_corporation_id"
	AllianceFactionID             string = "faction_id"
	AllianceIsClosed              string = "is_closed"
)

func NewAllianceRepository(db QueryExecContext) skillz.AllianceRepository {
	return &allianceRepository{
		db: db,
		alliance: tableConf{
			table: TableAlliances,
			columns: []string{
				AllianceID, AllianceName, AllianceTicker,
				AllianceDateFounded, AllianceCreatorID, AllianceCreatorCorporationID,
				AllianceExecutorCorporationID, AllianceFactionID, AllianceIsClosed,
				ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
	}
}

func (r *allianceRepository) Alliance(ctx context.Context, allianceID uint) (*skillz.Alliance, error) {

	query, args, err := sq.Select(r.alliance.columns...).
		From(r.alliance.table).
		Where(sq.Eq{CharacterID: allianceID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, allianceRepositoryIdentifier, "Alliance", "failed to generate sql")
	}

	var alliance = new(skillz.Alliance)
	err = r.db.GetContext(ctx, alliance, query, args...)
	return alliance, errors.Wrapf(err, prefixFormat, allianceRepositoryIdentifier, "Alliance")

}

func (r *allianceRepository) CreateAlliance(ctx context.Context, alliance *skillz.Alliance) error {

	now := time.Now()
	alliance.CreatedAt = now
	alliance.UpdatedAt = now

	query, args, err := sq.Insert(r.alliance.table).SetMap(map[string]interface{}{
		AllianceID:                    alliance.ID,
		AllianceName:                  alliance.Name,
		AllianceTicker:                alliance.Ticker,
		AllianceDateFounded:           alliance.DateFounded,
		AllianceCreatorID:             alliance.CreatorID,
		AllianceCreatorCorporationID:  alliance.CreatorCorporationID,
		AllianceExecutorCorporationID: alliance.ExecutorCorporationID,
		AllianceFactionID:             alliance.FactionID,
		AllianceIsClosed:              alliance.IsClosed,
		ColumnCreatedAt:               alliance.CreatedAt,
		ColumnUpdatedAt:               alliance.UpdatedAt,
	}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, allianceRepositoryIdentifier, "CreateAlliance", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, allianceRepositoryIdentifier, "CreateAlliance")

}

func (r *allianceRepository) UpdateAlliance(ctx context.Context, alliance *skillz.Alliance) error {
	alliance.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.alliance.table).SetMap(map[string]interface{}{
		AllianceName:                  alliance.Name,
		AllianceTicker:                alliance.Ticker,
		AllianceDateFounded:           alliance.DateFounded,
		AllianceCreatorID:             alliance.CreatorID,
		AllianceCreatorCorporationID:  alliance.CreatorCorporationID,
		AllianceExecutorCorporationID: alliance.ExecutorCorporationID,
		AllianceFactionID:             alliance.FactionID,
		AllianceIsClosed:              alliance.IsClosed,
		ColumnUpdatedAt:               alliance.UpdatedAt,
	}).Where(sq.Eq{AllianceID: alliance.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, allianceRepositoryIdentifier, "UpdateAlliance", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, allianceRepositoryIdentifier, "UpdateAlliance")
}
