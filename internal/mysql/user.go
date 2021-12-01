package mysql

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type userRepository struct {
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
	UserIsNew             = "is_new"
	UserDisabled          = "disabled"
	UserDisabledReason    = "disabled_reason"
	UserDisabledTimestamp = "disabled_timestamp"
	UserLastLogin         = "last_login"
)

func NewUserRepository(db QueryExecContext) skillz.UserRepository {
	return &userRepository{
		db: db,
		users: tableConf{
			table: TableUsers,
			columns: []string{
				UserID, UserCharacterID,
				UserAccessToken, UserRefreshToken,
				UserExpires, UserOwnerHash,
				UserScopes, UserIsNew, UserDisabled,
				UserDisabledReason, UserDisabledTimestamp,
				UserLastLogin, ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
	}
}

func (r *userRepository) User(ctx context.Context, id uuid.UUID) (*skillz.User, error) {

	query, args, err := sq.Select(r.users.columns...).
		From(r.users.table).
		Where(sq.Eq{UserID: id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, userRepositoryIdentifier, "User", "failed to generate sql")
	}

	var user = new(skillz.User)
	err = r.db.GetContext(ctx, user, query, args...)
	return user, errors.Wrapf(err, prefixFormat, userRepositoryIdentifier, "User")

}

func (r *userRepository) UserByCharacterID(ctx context.Context, characterID uint64) (*skillz.User, error) {

	query, args, err := sq.Select(r.users.columns...).
		From(r.users.table).
		Where(sq.Eq{UserCharacterID: characterID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, userRepositoryIdentifier, "UserByCharacterID", "failed to generate sql")
	}

	var user = new(skillz.User)
	err = r.db.GetContext(ctx, user, query, args...)
	return user, errors.Wrapf(err, prefixFormat, userRepositoryIdentifier, "UserByCharacterID")

}

func (r *userRepository) SearchUsers(ctx context.Context, q string) ([]*skillz.User, error) {

	query, args, err := sq.Select(r.users.columns...).
		From(r.users.table).
		InnerJoin(
			fmt.Sprintf(
				"%s on %s.%s = %s.%s",
				TableCharacters,
				TableCharacters,
				CharacterID,
				TableUsers,
				UserID,
			),
		).
		Where(sq.Like{fmt.Sprintf("%s.%s", TableCharacters, CharacterName): fmt.Sprintf("%%%s%%", q)}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, userRepositoryIdentifier, "SearchUsers", "failed to generate sql")
	}

	var users = make([]*skillz.User, 0)
	err = r.db.SelectContext(ctx, &users, query, args...)
	return users, errors.Wrapf(err, prefixFormat, userRepositoryIdentifier, "SearchUsers")

}

func (r *userRepository) CreateUser(ctx context.Context, user *skillz.User) error {

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
		return errors.Wrapf(err, errorFFormat, userRepositoryIdentifier, "CreateUser", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, userRepositoryIdentifier, "CreateUser")

}

func (r *userRepository) UpdateUser(ctx context.Context, user *skillz.User) error {

	user.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.users.table).SetMap(map[string]interface{}{
		UserCharacterID:       user.CharacterID,
		UserAccessToken:       user.AccessToken,
		UserRefreshToken:      user.RefreshToken,
		UserExpires:           user.Expires,
		UserOwnerHash:         user.OwnerHash,
		UserScopes:            user.Scopes,
		UserIsNew:             user.IsNew,
		UserDisabled:          user.Disabled,
		UserDisabledReason:    user.DisabledReason,
		UserDisabledTimestamp: user.DisabledTimestamp,
		UserLastLogin:         user.LastLogin,
		ColumnUpdatedAt:       user.UpdatedAt,
	}).Where(sq.Eq{UserID: user.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, userRepositoryIdentifier, "UpdateUser", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, userRepositoryIdentifier, "UpdateUser")

}
