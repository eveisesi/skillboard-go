package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type ETagRepository struct {
	db      QueryExecContext
	table   string
	columns []string
}

const (
	ETagPath        = "path"
	ETagETag        = "etag"
	ETagCachedUntil = "cached_until"
)

var _ skillz.EtagRepository = (*ETagRepository)(nil)

func NewETagRepository(db QueryExecContext) *ETagRepository {
	return &ETagRepository{
		db:    db,
		table: "etags",
		columns: []string{
			ETagPath, ETagETag, ETagCachedUntil, ColumnCreatedAt, ColumnUpdatedAt,
		},
	}
}

func (r *ETagRepository) Etag(ctx context.Context, path string) (*skillz.Etag, error) {

	query, args, err := sq.Select(r.columns...).
		From(r.table).
		Where(sq.Eq{ETagPath: path}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, etagRepository, "Etag", "failed to generate sql")
	}

	var etag = new(skillz.Etag)
	err = r.db.GetContext(ctx, etag, query, args...)
	return etag, errors.Wrapf(err, prefixFormat, etagRepository, "Etag")

}

func (r *ETagRepository) InsertEtag(ctx context.Context, etag *skillz.Etag) error {

	now := time.Now()
	etag.CreatedAt = now
	etag.UpdatedAt = now

	query, args, err := sq.Insert(r.table).SetMap(map[string]interface{}{
		ETagPath:        etag.Path,
		ETagETag:        etag.Etag,
		ETagCachedUntil: etag.CachedUntil,
		ColumnCreatedAt: etag.CreatedAt,
		ColumnUpdatedAt: etag.UpdatedAt,
	}).Suffix(OnDuplicateKeyStmt(
		ETagETag,
		ETagCachedUntil,
		ColumnUpdatedAt,
	)).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, etagRepository, "InsertEtag", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, etagRepository, "InsertEtag")

}
