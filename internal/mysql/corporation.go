package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type corporationRepository struct {
	db          QueryExecContext
	corporation tableConf
	history     tableConf
}

const (
	CorporationID            string = "id"
	CorporationAllianceID    string = "alliance_id"
	CorporationCeoID         string = "ceo_id"
	CorporationCreatorID     string = "creator_id"
	CorporationDateFounded   string = "date_founded"
	CorporationFactionID     string = "faction_id"
	CorporationHomeStationID string = "home_station_id"
	CorporationMemberCount   string = "member_count"
	CorporationName          string = "name"
	CorporationShares        string = "shares"
	CorporationTaxRate       string = "tax_rate"
	CorporationTicker        string = "ticker"
	CorporationURL           string = "url"
	CorporationWarEligible   string = "war_eligible"

	AllianceHistoryCorporationID string = "corporation_id"
	AllianceHistoryRecordID      string = "record_id"
	AllianceHistoryAllianceID    string = "alliance_id"
	AllianceHistoryIsDeleteed    string = "is_deleted"
	AllianceHistoryStartDate     string = "start_date"
)

func NewCorporationRepository(db QueryExecContext) skillz.CorporationRepository {

	return &corporationRepository{
		db: db,
		corporation: tableConf{
			table: TableCorporations,
			columns: []string{
				CorporationID, CorporationAllianceID, CorporationCeoID,
				CorporationCreatorID, CorporationDateFounded, CorporationFactionID,
				CorporationHomeStationID, CorporationMemberCount, CorporationName,
				CorporationShares, CorporationTaxRate, CorporationTicker,
				CorporationURL, CorporationWarEligible,
				ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
		history: tableConf{
			table: TableCorporationAllianceHistory,
			columns: []string{
				AllianceHistoryCorporationID, AllianceHistoryRecordID,
				AllianceHistoryAllianceID, AllianceHistoryIsDeleteed,
				AllianceHistoryStartDate, ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
	}
}

func (r *corporationRepository) Corporation(ctx context.Context, corporationID uint) (*skillz.Corporation, error) {

	query, args, err := sq.Select(r.corporation.columns...).
		From(r.corporation.table).
		Where(sq.Eq{CharacterID: corporationID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, corporationRepositoryIdentifier, "Corporation", "failed to generate sql")
	}

	var corporation = new(skillz.Corporation)
	err = r.db.GetContext(ctx, corporation, query, args...)
	return corporation, errors.Wrapf(err, prefixFormat, corporationRepositoryIdentifier, "Corporation")

}

func (r *corporationRepository) CreateCorporation(ctx context.Context, corporation *skillz.Corporation) error {

	now := time.Now()
	corporation.CreatedAt = now
	corporation.UpdatedAt = now

	query, args, err := sq.Insert(r.corporation.table).SetMap(map[string]interface{}{
		CorporationID:            corporation.ID,
		CorporationAllianceID:    corporation.AllianceID,
		CorporationCeoID:         corporation.CeoID,
		CorporationCreatorID:     corporation.CreatorID,
		CorporationDateFounded:   corporation.DateFounded,
		CorporationFactionID:     corporation.FactionID,
		CorporationHomeStationID: corporation.HomeStationID,
		CorporationMemberCount:   corporation.MemberCount,
		CorporationName:          corporation.Name,
		CorporationShares:        corporation.Shares,
		CorporationTaxRate:       corporation.TaxRate,
		CorporationTicker:        corporation.Ticker,
		CorporationURL:           corporation.URL,
		CorporationWarEligible:   corporation.WarEligible,
		ColumnCreatedAt:          corporation.CreatedAt,
		ColumnUpdatedAt:          corporation.UpdatedAt,
	}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, corporationRepositoryIdentifier, "CreateCorporation", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, corporationRepositoryIdentifier, "CreateCorporation")

}

func (r *corporationRepository) UpdateCorporation(ctx context.Context, corporation *skillz.Corporation) error {
	corporation.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.corporation.table).SetMap(map[string]interface{}{
		CorporationAllianceID:    corporation.AllianceID,
		CorporationCeoID:         corporation.CeoID,
		CorporationCreatorID:     corporation.CreatorID,
		CorporationDateFounded:   corporation.DateFounded,
		CorporationFactionID:     corporation.FactionID,
		CorporationHomeStationID: corporation.HomeStationID,
		CorporationMemberCount:   corporation.MemberCount,
		CorporationName:          corporation.Name,
		CorporationShares:        corporation.Shares,
		CorporationTaxRate:       corporation.TaxRate,
		CorporationTicker:        corporation.Ticker,
		CorporationURL:           corporation.URL,
		CorporationWarEligible:   corporation.WarEligible,
		ColumnUpdatedAt:          corporation.UpdatedAt,
	}).Where(sq.Eq{CorporationID: corporation.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, corporationRepositoryIdentifier, "UpdateCorporation", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, corporationRepositoryIdentifier, "UpdateCorporation")
}

func (r *corporationRepository) CorporationAllianceHistory(ctx context.Context, corporationID uint) ([]*skillz.CorporationAllianceHistory, error) {
	query, args, err := sq.Select(r.history.columns...).
		From(r.history.table).
		Where(sq.Eq{HistoryCorporationID: corporationID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, corporationRepositoryIdentifier, "CorporationAllianceHistory", "failed to generate sql")
	}

	var records = make([]*skillz.CorporationAllianceHistory, 0)
	err = r.db.SelectContext(ctx, &records, query, args...)
	return records, errors.Wrapf(err, prefixFormat, corporationRepositoryIdentifier, "CorporationAllianceHistory")
}

func (r *corporationRepository) CreateCorporationAllianceHistory(ctx context.Context, records []*skillz.CorporationAllianceHistory) ([]*skillz.CorporationAllianceHistory, error) {

	i := sq.Insert(r.history.table).Columns(r.history.columns...)
	now := time.Now()
	for _, record := range records {
		record.CreatedAt = now
		record.UpdatedAt = now
		i = i.Values(
			record.CorporationID,
			record.RecordID,
			record.AllianceID,
			record.IsDeleteed,
			record.StartDate,
			record.CreatedAt,
			record.UpdatedAt,
		)
	}

	i = i.Suffix(OnDuplicateKeyStmt(
		AllianceHistoryCorporationID,
		AllianceHistoryAllianceID,
		AllianceHistoryIsDeleteed,
		AllianceHistoryStartDate,
		ColumnUpdatedAt,
	))

	query, args, err := i.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, corporationRepositoryIdentifier, "CreateCorporationAllianceHistory", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return records, errors.Wrapf(err, prefixFormat, corporationRepositoryIdentifier, "CreateCorporationAllianceHistory")

}
