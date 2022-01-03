package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type cloneRepository struct {
	db             QueryExecContext
	meta, death    tableConf
	jump, implants tableConf
}

const (
	MetaLastCloneJumpDate     string = "last_clone_jump_date"
	MetaLastStationChangeDate string = "last_station_change_date"

	DeathLocationID   string = "location_id"
	DeathLocationType string = "location_type"

	JumpJumpCloneID  string = "jump_clone_id"
	JumpLocationID   string = "location_id"
	JumpLocationType string = "location_type"
	JumpImplantIDs   string = "implant_ids"

	ImplantsImplantID string = "implant_id"
)

func NewCloneRepository(db QueryExecContext) skillz.CloneRepository {
	return &cloneRepository{
		db: db,
		meta: tableConf{
			table: TableCharacterCloneMeta,
			columns: []string{
				MetaLastCloneJumpDate, MetaLastStationChangeDate,
				ColumnCharacterID, ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
		death: tableConf{
			table: TableCharacterHomeClone,
			columns: []string{
				DeathLocationID, DeathLocationType,
				ColumnCharacterID, ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
		jump: tableConf{
			table: TableCharacterJumpClones,
			columns: []string{
				JumpJumpCloneID, JumpLocationID, JumpLocationType, JumpImplantIDs,
				ColumnCharacterID, ColumnCreatedAt,
			},
		},
		implants: tableConf{
			table: TableCharacterImplants,
			columns: []string{
				ColumnCharacterID, ImplantsImplantID, ColumnCreatedAt,
			},
		},
	}
}

func (r *cloneRepository) CharacterCloneMeta(ctx context.Context, characterID uint64) (*skillz.CharacterCloneMeta, error) {

	query, args, err := sq.Select(r.meta.columns...).
		From(r.meta.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "CharacterCloneMeta", "failed to generate sql")
	}

	var meta = new(skillz.CharacterCloneMeta)
	err = r.db.GetContext(ctx, meta, query, args...)
	return meta, errors.Wrapf(err, prefixFormat, cloneRepositoryIdentifier, "CharacterCloneMeta")

}

func (r *cloneRepository) CreateCharacterCloneMeta(ctx context.Context, meta *skillz.CharacterCloneMeta) error {

	now := time.Now()
	meta.CreatedAt = now
	meta.UpdatedAt = now

	query, args, err := sq.Insert(r.meta.table).SetMap(map[string]interface{}{
		ColumnCharacterID:         meta.CharacterID,
		MetaLastCloneJumpDate:     meta.LastCloneJumpDate,
		MetaLastStationChangeDate: meta.LastStationChangeDate,
		ColumnCreatedAt:           meta.CreatedAt,
		ColumnUpdatedAt:           meta.UpdatedAt,
	}).
		Suffix(OnDuplicateKeyStmt(
			MetaLastCloneJumpDate,
			MetaLastStationChangeDate,
			ColumnUpdatedAt,
		)).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "CreateCloneMeta", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, cloneRepositoryIdentifier, "CreateCloneMeta")
}

func (r *cloneRepository) CharacterDeathClone(ctx context.Context, characterID uint64) (*skillz.CharacterDeathClone, error) {

	query, args, err := sq.Select(r.death.columns...).
		From(r.death.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "CharacterDeathClone", "failed to generate sql")
	}

	var death = new(skillz.CharacterDeathClone)
	err = r.db.GetContext(ctx, death, query, args...)
	return death, errors.Wrapf(err, prefixFormat, cloneRepositoryIdentifier, "CharacterDeathClone")

}

func (r *cloneRepository) CreateCharacterDeathClone(ctx context.Context, death *skillz.CharacterDeathClone) error {

	now := time.Now()
	death.CreatedAt = now
	death.UpdatedAt = now

	query, args, err := sq.Insert(r.death.table).SetMap(map[string]interface{}{
		ColumnCharacterID: death.CharacterID,
		DeathLocationID:   death.LocationID,
		DeathLocationType: death.LocationType,
		ColumnCreatedAt:   death.CreatedAt,
		ColumnUpdatedAt:   death.UpdatedAt,
	}).Suffix(OnDuplicateKeyStmt(
		DeathLocationID,
		DeathLocationType,
		ColumnUpdatedAt,
	)).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "CreateCharacterDeathClone", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, cloneRepositoryIdentifier, "CreateCharacterDeathClone")

}

func (r *cloneRepository) CharacterJumpClones(ctx context.Context, characterID uint64) ([]*skillz.CharacterJumpClone, error) {

	query, args, err := sq.Select(r.jump.columns...).
		From(r.jump.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "CharacterJumpClones", "failed to generate sql")
	}

	var jump = make([]*skillz.CharacterJumpClone, 0, 10)
	err = r.db.SelectContext(ctx, &jump, query, args...)
	return jump, errors.Wrapf(err, prefixFormat, cloneRepositoryIdentifier, "CharacterJumpClones")

}

func (r *cloneRepository) CreateCharacterJumpClones(ctx context.Context, clones []*skillz.CharacterJumpClone) error {

	i := sq.Insert(r.jump.table).Columns(r.jump.columns...)
	now := time.Now()
	for _, clone := range clones {
		clone.CreatedAt = now

		i = i.Values(
			clone.JumpCloneID, clone.LocationID,
			clone.LocationType, clone.ImplantIDs,
			clone.CharacterID, clone.CreatedAt,
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "CreateCharacterJumpClones", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *cloneRepository) DeleteCharacterJumpClones(ctx context.Context, characterID uint64) error {

	query, args, err := sq.Delete(r.jump.table).Where(sq.Eq{ColumnCharacterID: characterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "DeleteCharacterJumpClones", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *cloneRepository) CharacterImplants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error) {

	query, args, err := sq.Select(r.implants.columns...).
		From(r.implants.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
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
		i = i.Values(implant.CharacterID, implant.ImplantID, implant.CreatedAt)
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
