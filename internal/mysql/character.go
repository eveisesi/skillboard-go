package mysql

import (
	"context"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type tableConf struct {
	table   string
	columns []string
}

type CharacterRepository struct {
	db        *sqlx.DB
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

	HistoryCharacterID   string = "character_id"
	HistoryRecordID      string = "record_id"
	HistoryCorporationID string = "corporation_id"
	HistoryIsDeleted     string = "is_deleted"
	HistoryStartDate     string = "start_date"
)

var insertCharacterCorpHistoryDuplicateKeyStmt string

func NewCharacterRepository(db *sqlx.DB, character, history string) *CharacterRepository {

	t := make([]string, 0)
	for _, column := range []string{
		HistoryCharacterID, HistoryCorporationID,
		HistoryIsDeleted, HistoryStartDate, ColumnUpdatedAt,
	} {
		t = append(t, fmt.Sprintf("%[1]s = VALUES(%[1]s)", column))
	}
	insertCharacterCorpHistoryDuplicateKeyStmt = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", strings.Join(t, ","))

	return &CharacterRepository{
		character: tableConf{
			table: character,
			columns: []string{
				CharacterID, CharacterName, CharacterCorporationID,
				CharacterAllianceID, CharacterFactionID, CharacterSecurityStatus,
				CharacterGender, CharacterBirthday, CharacterTitle,
				CharacterBloodlineID, CharacterRaceID, ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
		history: tableConf{
			table: history,
			columns: []string{
				HistoryCharacterID, HistoryRecordID, HistoryCorporationID,
				HistoryIsDeleted, HistoryStartDate, ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
	}
}

func (r *CharacterRepository) Character(ctx context.Context, id uint64) (*skillz.Character, error) {

	query, args, err := sq.Select(r.character.columns...).
		From(r.character.table).
		Where(sq.Eq{CharacterID: id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate sql")
	}

	var character = new(skillz.Character)
	return character, r.db.GetContext(ctx, character, query, args...)
}

func (r *CharacterRepository) InsertCharacter(ctx context.Context, character *skillz.Character) (*skillz.Character, error) {

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
		return nil, errors.Wrap(err, "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args)
	return character, err
}

func (r *CharacterRepository) UpdateCharacter(ctx context.Context, character *skillz.Character) (*skillz.Character, error) {

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
		ColumnCreatedAt:         character.CreatedAt,
		ColumnUpdatedAt:         character.UpdatedAt,
	}).Where(sq.Eq{CharacterID: character.ID}).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args)
	return character, err

}

func (r *CharacterRepository) CharacterCorporationHistory(ctx context.Context, id uint64) ([]*skillz.CharacterCorporationHistory, error) {

	query, args, err := sq.Select(r.history.columns...).
		From(r.history.table).
		Where(sq.Eq{HistoryCharacterID: id}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate sql")
	}

	var records = make([]*skillz.CharacterCorporationHistory, 0, 256)
	return records, r.db.SelectContext(ctx, &records, query, args...)

}

func (r *CharacterRepository) CreateCharacterCorporationHistory(ctx context.Context, records []*skillz.CharacterCorporationHistory) ([]*skillz.CharacterCorporationHistory, error) {

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

	i = i.Suffix(insertCharacterCorpHistoryDuplicateKeyStmt)

	query, args, err := i.ToSql()
	if err != nil {
		return nil, err
	}

	fmt.Println(query)

	_, err = r.db.ExecContext(ctx, query, args...)
	return records, err

}
