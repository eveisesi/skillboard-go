package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type skillRepository struct {
	db               QueryExecContext
	attributes, meta tableConf
	skills, queue    tableConf
}

const (
	MetaTotalSP       string = "total_sp"
	MetaUnallocatedSP string = "unallocated_sp"

	SkillsActiveSkillLevel   string = "active_skill_level"
	SkillsSkillID            string = "skill_id"
	SkillsSkillpointsInSkill string = "skillpoints_in_skill"
	SkillsTrainedSkillLevel  string = "trained_skill_level"

	QueueQueuePosition   string = "queue_position"
	QueueSkillID         string = "skill_id"
	QueueFinishedLevel   string = "finished_level"
	QueueTrainingStartSP string = "training_start_sp"
	QueueLevelStartSP    string = "level_start_sp"
	QueueLevelEndSP      string = "level_end_sp"
	QueueStartDate       string = "start_date"
	QueueFinishDate      string = "finish_date"

	AttributesCharisma                 string = "charisma"
	AttributesIntelligence             string = "intelligence"
	AttributesMemory                   string = "memory"
	AttributesPerception               string = "perception"
	AttributesWillpower                string = "willpower"
	AttributesBonusRemaps              string = "bonus_remaps"
	AttributesLastRemapDate            string = "last_remap_date"
	AttributesAccruedRemapCooldownDate string = "accrued_remap_cooldown_date"
)

func NewSkillRepository(db QueryExecContext) skillz.CharacterSkillRepository {

	return &skillRepository{
		db: db,
		attributes: tableConf{
			table: TableCharacterAttributes,
			columns: []string{
				AttributesCharisma, AttributesIntelligence,
				AttributesMemory, AttributesPerception,
				AttributesWillpower, AttributesBonusRemaps,
				AttributesLastRemapDate, AttributesAccruedRemapCooldownDate,
				ColumnCharacterID, ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
		meta: tableConf{
			table: TableCharacterSkillMeta,
			columns: []string{
				MetaTotalSP, MetaUnallocatedSP,
				ColumnCharacterID, ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
		skills: tableConf{
			table: TableCharacterSkills,
			columns: []string{
				SkillsActiveSkillLevel, SkillsSkillID,
				SkillsSkillpointsInSkill, SkillsTrainedSkillLevel,
				ColumnCharacterID, ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
		queue: tableConf{
			table: TableCharacterSkillQueue,
			columns: []string{
				QueueQueuePosition, QueueSkillID,
				QueueFinishedLevel, QueueTrainingStartSP,
				QueueLevelStartSP, QueueLevelEndSP,
				QueueStartDate, QueueFinishDate,
				ColumnCharacterID, ColumnCreatedAt,
			},
		},
	}

}

func (r *SkillRepository) CharacterAttributes(ctx context.Context, characterID uint64) (*skillz.CharacterAttributes, error) {

	query, args, err := sq.Select(r.attributes.columns...).
		From(r.attributes.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "CharacterAttributes", "failed to generate sql")
	}

	var attributes = new(skillz.CharacterAttributes)
	err = r.db.GetContext(ctx, attributes, query, args...)
	return attributes, errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "CharacterAttributes")

}

func (r *SkillRepository) CreateCharacterAttributes(ctx context.Context, attributes *skillz.CharacterAttributes) error {

	now := time.Now()
	attributes.CreatedAt = now
	attributes.UpdatedAt = now

	query, args, err := sq.Insert(r.attributes.table).SetMap(map[string]interface{}{
		AttributesCharisma:                 attributes.Charisma,
		AttributesIntelligence:             attributes.Intelligence,
		AttributesMemory:                   attributes.Memory,
		AttributesPerception:               attributes.Perception,
		AttributesWillpower:                attributes.Willpower,
		AttributesBonusRemaps:              attributes.BonusRemaps,
		AttributesLastRemapDate:            attributes.LastRemapDate,
		AttributesAccruedRemapCooldownDate: attributes.AccruedRemapCooldownDate,
		ColumnCharacterID:                  attributes.CharacterID,
		ColumnCreatedAt:                    attributes.CreatedAt,
		ColumnUpdatedAt:                    attributes.UpdatedAt,
	}).
		Suffix(OnDuplicateKeyStmt(
			AttributesCharisma,
			AttributesIntelligence,
			AttributesMemory,
			AttributesPerception,
			AttributesWillpower,
			AttributesBonusRemaps,
			AttributesLastRemapDate,
			AttributesAccruedRemapCooldownDate,
			ColumnUpdatedAt,
		)).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "CreateCharacterAttributes", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "CreateCharacterAttributes")

}

func (r *SkillRepository) UpdateCharacterAttributes(ctx context.Context, attributes *skillz.CharacterAttributes) error {

	attributes.UpdatedAt = time.Now()
	query, args, err := sq.Update(r.attributes.table).SetMap(map[string]interface{}{
		AttributesCharisma:                 attributes.Charisma,
		AttributesIntelligence:             attributes.Intelligence,
		AttributesMemory:                   attributes.Memory,
		AttributesPerception:               attributes.Perception,
		AttributesWillpower:                attributes.Willpower,
		AttributesBonusRemaps:              attributes.BonusRemaps,
		AttributesLastRemapDate:            attributes.LastRemapDate,
		AttributesAccruedRemapCooldownDate: attributes.AccruedRemapCooldownDate,
		ColumnUpdatedAt:                    attributes.UpdatedAt,
	}).Where(sq.Eq{ColumnCharacterID: attributes.CharacterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "UpdateCharacterAttributes", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "UpdateCharacterAttributes")

}

func (r *SkillRepository) DeleteCharacterAttributes(ctx context.Context, characterID uint64) error {

	query, args, err := sq.Delete(r.attributes.table).Where(sq.Eq{ColumnCharacterID: characterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "DeleteCharacterAttributes", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "DeleteCharacterAttributes")

}

func (r *SkillRepository) CharacterSkillMeta(ctx context.Context, characterID uint64) (*skillz.CharacterSkillMeta, error) {

	query, args, err := sq.Select(r.meta.columns...).
		From(r.meta.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "CharacterSkillMeta", "failed to generate sql")
	}

	var meta = new(skillz.CharacterSkillMeta)
	err = r.db.GetContext(ctx, meta, query, args...)
	return meta, errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "CharacterSkillMeta")

}

func (r *SkillRepository) CreateCharacterSkillMeta(ctx context.Context, meta *skillz.CharacterSkillMeta) error {

	now := time.Now()
	meta.CreatedAt = now
	meta.UpdatedAt = now

	query, args, err := sq.Insert(r.meta.table).SetMap(map[string]interface{}{
		MetaTotalSP:       meta.TotalSP,
		MetaUnallocatedSP: meta.UnallocatedSP,
		ColumnCharacterID: meta.CharacterID,
		ColumnCreatedAt:   meta.CreatedAt,
		ColumnUpdatedAt:   meta.UpdatedAt,
	}).
		Suffix(OnDuplicateKeyStmt(MetaTotalSP, MetaUnallocatedSP, ColumnUpdatedAt)).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "CreateCharacterSkillMeta", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "CreateCharacterSkillMeta")

}

// func (r *SkillRepository) UpdateCharacterSkillMeta(ctx context.Context, meta *skillz.CharacterSkillMeta) error {

// 	meta.UpdatedAt = time.Now()

// 	query, args, err := sq.Update(r.meta.table).SetMap(map[string]interface{}{
// 		MetaTotalSP:       meta.TotalSP,
// 		MetaUnallocatedSP: meta.UnallocatedSP,

// 		ColumnUpdatedAt: meta.UpdatedAt,
// 	}).Where(sq.Eq{ColumnCharacterID: meta.CharacterID}).ToSql()
// 	if err != nil {
// 		return errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "UpdateCharacterSkillMeta", "failed to generate sql")
// 	}

// 	_, err = r.db.ExecContext(ctx, query, args...)
// 	return errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "UpdateCharacterSkillMeta")

// }

func (r *SkillRepository) DeleteCharacterSkillMeta(ctx context.Context, characterID uint64) error {

	query, args, err := sq.Delete(r.meta.table).Where(sq.Eq{ColumnCharacterID: characterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "DeleteCharacterSkillMeta", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "DeleteCharacterSkillMeta")

}

func (r *SkillRepository) CharacterSkills(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkill, error) {

	query, args, err := sq.Select(r.skills.columns...).
		From(r.skills.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "CharacterSkills", "failed to generate sql")
	}

	var skills = make([]*skillz.CharacterSkill, 0, 1024)
	err = r.db.SelectContext(ctx, &skills, query, args...)
	return skills, errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "CharacterSkills")

}

func (r *SkillRepository) CreateCharacterSkills(ctx context.Context, skills []*skillz.CharacterSkill) error {

	now := time.Now()

	i := sq.Insert(r.skills.table).Columns(r.skills.columns...)
	for _, skill := range skills {
		skill.CreatedAt = now
		skill.UpdatedAt = now

		i = i.Values(
			skill.ActiveSkillLevel,
			skill.SkillID,
			skill.SkillpointsInSkill,
			skill.TrainedSkillLevel,
			skill.CharacterID,
			skill.CreatedAt,
			skill.UpdatedAt,
		)
	}
	i = i.Suffix(OnDuplicateKeyStmt(SkillsActiveSkillLevel, SkillsSkillpointsInSkill, SkillsTrainedSkillLevel, ColumnUpdatedAt))

	query, args, err := i.ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "CreateCharacterSkillMeta", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "CreateCharacterSkillMeta")

}

func (r *SkillRepository) DeleteCharacterSkills(ctx context.Context, characterID uint64) error {

	query, args, err := sq.Delete(r.skills.table).Where(sq.Eq{ColumnCharacterID: characterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "DeleteCharacterSkills", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "DeleteCharacterSkills")

}

func (r *SkillRepository) CharacterSkillQueue(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillQueue, error) {

	query, args, err := sq.Select(r.queue.columns...).
		From(r.queue.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "CharacterSkillQueue", "failed to generate sql")
	}

	var positions = make([]*skillz.CharacterSkillQueue, 0, 50)
	err = r.db.SelectContext(ctx, &positions, query, args...)
	return positions, errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "CharacterSkillQueue")

}

func (r *SkillRepository) CreateCharacterSkillQueue(ctx context.Context, positions []*skillz.CharacterSkillQueue) error {

	now := time.Now()

	i := sq.Insert(r.queue.table).Columns(r.queue.columns...)
	for _, position := range positions {
		position.CreatedAt = now

		i = i.Values(
			position.QueuePosition, position.SkillID,
			position.FinishedLevel, position.TrainingStartSp,
			position.LevelStartSp, position.LevelEndSp,
			position.StartDate, position.FinishDate,
			position.CharacterID, position.CreatedAt,
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "CreateCharacterSkillQueue", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "CreateCharacterSkillQueue")

}

func (r *SkillRepository) DeleteCharacterSkillQueue(ctx context.Context, characterID uint64) error {

	query, args, err := sq.Delete(r.queue.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillsRepositoryIdentifier, "DeleteCharacterSkillQueue", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, skillsRepositoryIdentifier, "DeleteCharacterSkillQueue")

}
