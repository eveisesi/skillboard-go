package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type characterRepository struct {
	db        QueryExecContext
	character tableConf
	history   tableConf
}

const (
	CharacterID             string = "id"
	CharacterName           string = "name"
	CharacterCorporationID  string = "corporation_id"
	CharacterAllianceID     string = "alliance_id"
	CharacterFactionID      string = "faction_id"
	CharacterSecurityStatus string = "security_status"
	CharacterGender         string = "gender"
	CharacterBirthday       string = "birthday"
	CharacterTitle          string = "title"
	CharacterBloodlineID    string = "bloodline_id"
	CharacterRaceID         string = "race_id"

	HistoryRecordID      string = "record_id"
	HistoryCorporationID string = "corporation_id"
	HistoryIsDeleted     string = "is_deleted"
	HistoryStartDate     string = "start_date"
)

func NewCharacterRepository(db QueryExecContext) skillz.CharacterRepository {

	return &characterRepository{
		db: db,
		character: tableConf{
			table: TableCharacters,
			columns: []string{
				CharacterID, CharacterName, CharacterCorporationID,
				CharacterAllianceID, CharacterFactionID, CharacterSecurityStatus,
				CharacterGender, CharacterBirthday, CharacterTitle,
				CharacterBloodlineID, CharacterRaceID, ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
		history: tableConf{
			table: TableCharacterCorporationHistory,
			columns: []string{
				ColumnCharacterID, HistoryRecordID, HistoryCorporationID,
				HistoryIsDeleted, HistoryStartDate, ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
	}
}

func (r *characterRepository) Character(ctx context.Context, characterID uint64) (*skillz.Character, error) {

	query, args, err := sq.Select(r.character.columns...).
		From(r.character.table).
		Where(sq.Eq{CharacterID: characterID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, characterRepositoryIdentifier, "Character", "failed to generate sql")
	}

	var character = new(skillz.Character)
	err = r.db.GetContext(ctx, character, query, args...)
	return character, errors.Wrapf(err, prefixFormat, characterRepositoryIdentifier, "Character")

}

func (r *characterRepository) CreateCharacter(ctx context.Context, character *skillz.Character) error {

	now := time.Now()
	character.CreatedAt = now
	character.UpdatedAt = now

	query, args, err := sq.Insert(r.character.table).SetMap(map[string]interface{}{
		CharacterID:             character.ID,
		CharacterName:           character.Name,
		CharacterCorporationID:  character.CorporationID,
		CharacterAllianceID:     character.AllianceID,
		CharacterFactionID:      character.FactionID,
		CharacterSecurityStatus: character.SecurityStatus,
		CharacterGender:         character.Gender,
		CharacterBirthday:       character.Birthday,
		CharacterTitle:          character.Title,
		CharacterBloodlineID:    character.BloodlineID,
		CharacterRaceID:         character.RaceID,
		ColumnCreatedAt:         character.CreatedAt,
		ColumnUpdatedAt:         character.UpdatedAt,
	}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, characterRepositoryIdentifier, "CreateCharacter", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, characterRepositoryIdentifier, "CreateCharacter")
}

func (r *characterRepository) UpdateCharacter(ctx context.Context, character *skillz.Character) error {

	character.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.character.table).SetMap(map[string]interface{}{
		CharacterName:           character.Name,
		CharacterCorporationID:  character.CorporationID,
		CharacterAllianceID:     character.AllianceID,
		CharacterFactionID:      character.FactionID,
		CharacterSecurityStatus: character.SecurityStatus,
		CharacterGender:         character.Gender,
		CharacterBirthday:       character.Birthday,
		CharacterTitle:          character.Title,
		CharacterBloodlineID:    character.BloodlineID,
		CharacterRaceID:         character.RaceID,
		ColumnUpdatedAt:         character.UpdatedAt,
	}).Where(sq.Eq{CharacterID: character.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, characterRepositoryIdentifier, "UpdateCharacter", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, characterRepositoryIdentifier, "UpdateCharacter")

}

func (r *characterRepository) CharacterCorporationHistory(ctx context.Context, characterID uint64) ([]*skillz.CharacterCorporationHistory, error) {

	query, args, err := sq.Select(r.history.columns...).
		From(r.history.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, characterRepositoryIdentifier, "CharacterCorporationHistory", "failed to generate sql")
	}

	var records = make([]*skillz.CharacterCorporationHistory, 0, 256)
	err = r.db.SelectContext(ctx, &records, query, args...)
	return records, errors.Wrapf(err, prefixFormat, characterRepositoryIdentifier, "CharacterCorporationHistory")

}

func (r *characterRepository) CreateCharacterCorporationHistory(ctx context.Context, records []*skillz.CharacterCorporationHistory) ([]*skillz.CharacterCorporationHistory, error) {

	i := sq.Insert(r.history.table).Columns(r.history.columns...)
	now := time.Now()
	for _, record := range records {
		record.CreatedAt = now
		record.UpdatedAt = now
		i = i.Values(
			record.CharacterID,
			record.RecordID,
			record.CorporationID,
			record.IsDeleted,
			record.StartDate,
			record.CreatedAt,
			record.UpdatedAt,
		)
	}

	i = i.Suffix(OnDuplicateKeyStmt(
		ColumnCharacterID,
		HistoryCorporationID,
		HistoryIsDeleted,
		HistoryStartDate,
		ColumnUpdatedAt,
	))

	query, args, err := i.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, characterRepositoryIdentifier, "CreateCharacterCorporationHistory", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return records, errors.Wrapf(err, prefixFormat, characterRepositoryIdentifier, "CreateCharacterCorporationHistory")

}
