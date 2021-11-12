package skillz

import (
	"context"
	"time"
)

type EtagRepository interface {
	Etag(ctx context.Context, path string) (*Etag, error)
	// Etags(ctx context.Context, operators ...*Operator) ([]*Etag, error)
	InsertEtag(ctx context.Context, etag *Etag) error
	// UpdateEtag(ctx context.Context, path string, etag *Etag) (*Etag, error)
	// DeleteEtag(ctx context.Context, path string) (bool, error)
}

type Etag struct {
	Path        string    `db:"path" json:"path"`
	Etag        string    `db:"etag" json:"etag"`
	CachedUntil time.Time `db:"cached_until" json:"cached_until"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
