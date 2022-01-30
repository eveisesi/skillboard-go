package mysql

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type cloneRepository struct {
	db       QueryExecContext
	implants tableConf
}

const (
	ImplantsImplantID string = "implant_id"
	ImplantsSlot      string = "slot"
)

func NewCloneRepository(db QueryExecContext) skillz.CloneRepository {
	return &cloneRepository{
		db: db,
		implants: tableConf{
			table: TableCharacterImplants,
			columns: []string{
				ColumnCharacterID, ImplantsImplantID, ImplantsSlot, ColumnCreatedAt,
			},
		},
	}
}

func (r *cloneRepository) CharacterImplants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error) {

	query, args, err := sq.Select(r.implants.columns...).
		From(r.implants.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		Where(sq.LtOrEq{ImplantsSlot: 5}).
		OrderBy(fmt.Sprintf("%s %s", ImplantsSlot, "ASC")).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "CharacterImplants", "failed to generate sql")
	}

	var implants = make([]*skillz.CharacterImplant, 0, 10)
	err = r.db.SelectContext(ctx, &implants, query, args...)
	return implants, errors.Wrapf(err, prefixFormat, cloneRepositoryIdentifier, "CharacterImplants")

}

func (r *cloneRepository) CreateCharacterImplants(ctx context.Context, implants []*skillz.CharacterImplant) error {

	i := sq.Insert(r.implants.table).Columns(r.implants.columns...)
	now := time.Now()
	for _, implant := range implants {
		implant.CreatedAt = now
		i = i.Values(implant.CharacterID, implant.ImplantID, implant.Slot, implant.CreatedAt)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "CreateCharacterImplants", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *cloneRepository) DeleteCharacterImplants(ctx context.Context, characterID uint64) error {

	query, args, err := sq.Delete(r.implants.table).Where(sq.Eq{ColumnCharacterID: characterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "DeleteCharacterImplants", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}
