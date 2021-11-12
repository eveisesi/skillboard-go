package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type CloneRepository struct {
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
	JumpImplants     string = "implants"

	ImplantsImplantID string = "implant_id"
)

var _ skillz.CloneRepository = new(CloneRepository)

func NewCloneRepository(db QueryExecContext, meta, death, jump, implants string) *CloneRepository {
	return &CloneRepository{
		db: db,
		meta: tableConf{
			table: meta,
			columns: []string{
				MetaLastCloneJumpDate, MetaLastStationChangeDate,
				ColumnCharacterID, ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
		death: tableConf{
			table: death,
			columns: []string{
				DeathLocationID, DeathLocationType,
				ColumnCharacterID, ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
		jump: tableConf{
			table: jump,
			columns: []string{
				JumpJumpCloneID, JumpLocationID, JumpLocationType, JumpImplants,
				ColumnCharacterID, ColumnCreatedAt,
			},
		},
		implants: tableConf{
			table: implants,
			columns: []string{
				ImplantsImplantID,
				ColumnCharacterID, ColumnCreatedAt,
			},
		},
	}
}

func (r *CloneRepository) CharacterCloneMeta(ctx context.Context, characterID uint64) (*skillz.CharacterCloneMeta, error) {

	query, args, err := sq.Select(r.meta.columns...).
		From(r.meta.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, cloneRepository, "CharacterCloneMeta", "failed to generate sql")
	}

	var meta = new(skillz.CharacterCloneMeta)
	err = r.db.GetContext(ctx, meta, query, args...)
	return meta, errors.Wrapf(err, prefixFormat, cloneRepository, "CharacterCloneMeta")

}

func (r *CloneRepository) CreateCloneMeta(ctx context.Context, meta *skillz.CharacterCloneMeta) error {

	now := time.Now()
	meta.CreatedAt = now
	meta.UpdatedAt = now

	query, args, err := sq.Insert(r.meta.table).SetMap(map[string]interface{}{
		ColumnCharacterID:         meta.CharacterID,
		MetaLastCloneJumpDate:     meta.LastCloneJumpDate,
		MetaLastStationChangeDate: meta.LastStationChangeDate,
		ColumnCreatedAt:           meta.CreatedAt,
		ColumnUpdatedAt:           meta.UpdatedAt,
	}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepository, "CreateCloneMeta", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, cloneRepository, "CreateCloneMeta")
}

func (r *CloneRepository) UpdateCloneMeta(ctx context.Context, meta *skillz.CharacterCloneMeta) error {

	meta.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.meta.table).SetMap(map[string]interface{}{
		MetaLastCloneJumpDate:     meta.LastCloneJumpDate,
		MetaLastStationChangeDate: meta.LastStationChangeDate,
		ColumnUpdatedAt:           meta.UpdatedAt,
	}).Where(sq.Eq{ColumnCharacterID: meta.CharacterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepository, "UpdateCloneMeta", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, cloneRepository, "UpdateCloneMeta")
}

func (r *CloneRepository) CharacterDeathClone(ctx context.Context, characterID uint64) (*skillz.CharacterDeathClone, error) {

	query, args, err := sq.Select(r.death.columns...).
		From(r.death.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, cloneRepository, "CharacterDeathClone", "failed to generate sql")
	}

	var death = new(skillz.CharacterDeathClone)
	err = r.db.GetContext(ctx, death, query, args...)
	return death, errors.Wrapf(err, prefixFormat, cloneRepository, "CharacterDeathClone")

}

func (r *CloneRepository) CreateCharacterDeathClone(ctx context.Context, death *skillz.CharacterDeathClone) error {

	now := time.Now()
	death.CreatedAt = now
	death.UpdatedAt = now

	query, args, err := sq.Insert(r.death.table).SetMap(map[string]interface{}{
		ColumnCharacterID: death.CharacterID,
		DeathLocationID:   death.LocationID,
		DeathLocationType: death.LocationType,
		ColumnCreatedAt:   death.CreatedAt,
		ColumnUpdatedAt:   death.UpdatedAt,
	}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepository, "CreateCharacterDeathClone", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, cloneRepository, "CreateCharacterDeathClone")

}

func (r *CloneRepository) UpdateCharacterDeathClone(ctx context.Context, death *skillz.CharacterDeathClone) error {

	death.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.death.table).SetMap(map[string]interface{}{
		DeathLocationID:   death.LocationID,
		DeathLocationType: death.LocationType,
		ColumnUpdatedAt:   death.UpdatedAt,
	}).Where(sq.Eq{ColumnCharacterID: death.CharacterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepository, "CreateCharacterDeathClone", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, cloneRepository, "CreateCharacterDeathClone")

}

func (r *CloneRepository) CharacterJumpClones(ctx context.Context, characterID uint64) ([]*skillz.CharacterJumpClone, error) {

	query, args, err := sq.Select(r.jump.columns...).
		From(r.jump.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, cloneRepository, "CharacterJumpClones", "failed to generate sql")
	}

	var jump = make([]*skillz.CharacterJumpClone, 0, 10)
	err = r.db.SelectContext(ctx, &jump, query, args...)
	return jump, errors.Wrapf(err, prefixFormat, cloneRepository, "CharacterJumpClones")

}

func (r *CloneRepository) CreateCharacterJumpClones(ctx context.Context, clones []*skillz.CharacterJumpClone) error {

	i := sq.Insert(r.jump.table).Columns(r.jump.columns...)
	now := time.Now()
	for _, clone := range clones {
		clone.CreatedAt = now

		i = i.Values(
			clone.JumpCloneID, clone.LocationID,
			clone.LocationType, clone.Implants,
			clone.CharacterID, clone.CreatedAt,
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepository, "CreateCharacterJumpClones", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *CloneRepository) DeleteCharacterJumpClones(ctx context.Context, characterID uint64) error {

	query, args, err := sq.Delete(r.jump.table).Where(sq.Eq{ColumnCharacterID: characterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepository, "DeleteCharacterJumpClones", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *CloneRepository) CharacterImplants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error) {

	query, args, err := sq.Select(r.implants.columns...).
		From(r.implants.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, cloneRepository, "CharacterImplants", "failed to generate sql")
	}

	var implants = make([]*skillz.CharacterImplant, 0, 10)
	err = r.db.SelectContext(ctx, &implants, query, args...)
	return implants, errors.Wrapf(err, prefixFormat, cloneRepository, "CharacterImplants")

}

func (r *CloneRepository) CreateCharacterImplants(ctx context.Context, implants []*skillz.CharacterImplant) error {

	i := sq.Insert(r.implants.table).Columns(r.implants.columns...)
	for _, implant := range implants {
		i.Values(implant.CharacterID, implant.ImplantID, implant.CreatedAt)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepository, "CreateCharacterImplants", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *CloneRepository) DeleteCharacterImplants(ctx context.Context, characterID uint64) error {

	query, args, err := sq.Delete(r.implants.table).Where(sq.Eq{ColumnCharacterID: characterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepository, "DeleteCharacterImplants", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}
