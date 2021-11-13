package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type UserRepository struct {
	db    QueryExecContext
	users tableConf
}

const (
	UserID                = "id"
	UserCharacterID       = "character_id"
	UserAccessToken       = "access_token"
	UserRefreshToken      = "refresh_token"
	UserExpires           = "expires"
	UserOwnerHash         = "owner_hash"
	UserScopes            = "scopes"
	UserDisabled          = "disabled"
	UserDisabledReason    = "disabled_reason"
	UserDisabledTimestamp = "disabled_timestamp"
	UserLastLogin         = "last_login"
)

var _ skillz.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db QueryExecContext) *UserRepository {
	return &UserRepository{
		db: db,
		users: tableConf{
			table: "users",
			columns: []string{
				UserID, UserCharacterID,
				UserAccessToken, UserRefreshToken,
				UserExpires, UserOwnerHash,
				UserScopes, UserDisabled,
				UserDisabledReason, UserDisabledTimestamp,
				UserLastLogin, ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
	}
}

func (r *UserRepository) User(ctx context.Context, id uuid.UUID) (*skillz.User, error) {

	query, args, err := sq.Select(r.users.columns...).
		From(r.users.table).
		Where(sq.Eq{UserID: id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, userRepository, "User", "failed to generate sql")
	}

	var user = new(skillz.User)
	err = r.db.GetContext(ctx, user, query, args...)
	return user, errors.Wrapf(err, prefixFormat, userRepository, "User")

}

func (r *UserRepository) UserByCharacterID(ctx context.Context, characterID uint64) (*skillz.User, error) {

	query, args, err := sq.Select(r.users.columns...).
		From(r.users.table).
		Where(sq.Eq{UserCharacterID: characterID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, userRepository, "UserByCharacterID", "failed to generate sql")
	}

	var user = new(skillz.User)
	err = r.db.GetContext(ctx, user, query, args...)
	return user, errors.Wrapf(err, prefixFormat, userRepository, "UserByCharacterID")

}

func (r *UserRepository) CreateUser(ctx context.Context, user *skillz.User) error {

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query, args, err := sq.Insert(r.users.table).SetMap(map[string]interface{}{
		UserID:                user.ID,
		UserCharacterID:       user.CharacterID,
		UserAccessToken:       user.AccessToken,
		UserRefreshToken:      user.RefreshToken,
		UserExpires:           user.Expires,
		UserOwnerHash:         user.OwnerHash,
		UserScopes:            user.Scopes,
		UserDisabled:          user.Disabled,
		UserDisabledReason:    user.DisabledReason,
		UserDisabledTimestamp: user.DisabledTimestamp,
		UserLastLogin:         user.LastLogin,
		ColumnCreatedAt:       user.CreatedAt,
		ColumnUpdatedAt:       user.UpdatedAt,
	}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, userRepository, "CreateUser", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, userRepository, "CreateUser")

}

func (r *UserRepository) UpdateUser(ctx context.Context, user *skillz.User) error {

	user.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.users.table).SetMap(map[string]interface{}{
		UserCharacterID:       user.CharacterID,
		UserAccessToken:       user.AccessToken,
		UserRefreshToken:      user.RefreshToken,
		UserExpires:           user.Expires,
		UserOwnerHash:         user.OwnerHash,
		UserScopes:            user.Scopes,
		UserDisabled:          user.Disabled,
		UserDisabledReason:    user.DisabledReason,
		UserDisabledTimestamp: user.DisabledTimestamp,
		UserLastLogin:         user.LastLogin,
		ColumnUpdatedAt:       user.UpdatedAt,
	}).Where(sq.Eq{UserID: user.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, userRepository, "UpdateUser", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, userRepository, "UpdateUser")

}
