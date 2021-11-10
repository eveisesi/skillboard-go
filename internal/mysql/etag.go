package mysql

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ETagRepository struct {
	db      *sqlx.DB
	table   string
	columns []string
}

const (
	ETagPath        = "path"
	ETagETag        = "etag"
	ETagCachedUntil = "cached_until"
)

func NewETagRepository(db *sqlx.DB, table string) *ETagRepository {
	return &ETagRepository{
		db:    db,
		table: table,
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
		return nil, errors.Wrap(err, "failed to generate sql")
	}

	var etag = new(skillz.Etag)
	return etag, r.db.GetContext(ctx, etag, query, args...)

}

var insertEtagDuplicateKeyStmt = fmt.Sprintf(
	"ON DUPLICATE KEY UPDATE %[1]s = VALUES(%[1]s), %[2]s = VALUES(%[2]s), %[3]s = VALLUES(%[3]s)",
	ETagETag,
	ETagCachedUntil,
	ColumnUpdatedAt,
)

func (r *ETagRepository) InsertEtag(ctx context.Context, etag *skillz.Etag) (*skillz.Etag, error) {

	now := time.Now()
	etag.CreatedAt = now
	etag.UpdatedAt = now

	query, args, err := sq.Insert(r.table).SetMap(map[string]interface{}{
		ETagPath:        etag.Path,
		ETagETag:        etag.Etag,
		ETagCachedUntil: etag.CachedUntil,
		ColumnCreatedAt: etag.CreatedAt,
		ColumnUpdatedAt: etag.UpdatedAt,
	}).Suffix(insertEtagDuplicateKeyStmt).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate sql")
	}

	fmt.Println(query)

	_, err = r.db.ExecContext(ctx, query, args)
	return etag, err

}
