package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type etagRepository struct {
	db      QueryExecContext
	table   string
	columns []string
}

const (
	ETagPath        = "path"
	ETagETag        = "etag"
	ETagCachedUntil = "cached_until"
)

func NewETagRepository(db QueryExecContext) skillz.EtagRepository {
	return &etagRepository{
		db:    db,
		table: TableEtags,
		columns: []string{
			ETagPath, ETagETag, ETagCachedUntil, ColumnCreatedAt, ColumnUpdatedAt,
		},
	}
}

func (r *etagRepository) Etag(ctx context.Context, path string) (*skillz.Etag, error) {

	query, args, err := sq.Select(r.columns...).
		From(r.table).
		Where(sq.Eq{ETagPath: path}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, etagRepositoryIdentifier, "Etag", "failed to generate sql")
	}

	var etag = new(skillz.Etag)
	err = r.db.GetContext(ctx, etag, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, prefixFormat, etagRepositoryIdentifier, "Etag")
	}
	return etag, nil

}

func (r *etagRepository) InsertEtag(ctx context.Context, etag *skillz.Etag) error {

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
		return errors.Wrapf(err, errorFFormat, etagRepositoryIdentifier, "InsertEtag", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, etagRepositoryIdentifier, "InsertEtag")

}
