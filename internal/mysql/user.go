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
	db       QueryExecContext
	users    tableConf
	settings tableConf
}

const (
	UserID                = "id"
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
	UserLastProcessed     = "last_processed"
	SettingsUserID        = "user_id"
	SettingsHideClones    = "hide_clones"
	SettingsHideQueue     = "hide_queue"
	SettingsHideStandings = "hide_standings"
	SettingsHideShips     = "hide_ships"
)

func NewUserRepository(db QueryExecContext) skillz.UserRepository {
	return &userRepository{
		db: db,
		users: tableConf{
			table: TableUsers,
			columns: []string{
				UserID, ColumnCharacterID,
				UserAccessToken, UserRefreshToken,
				UserExpires, UserOwnerHash,
				UserScopes, UserIsNew, UserDisabled,
				UserDisabledReason, UserDisabledTimestamp,
				UserLastLogin, ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
		settings: tableConf{
			table: TableUserSettings,
			columns: []string{
				SettingsUserID, SettingsHideClones,
				SettingsHideQueue, SettingsHideStandings, SettingsHideShips,
				ColumnCreatedAt, ColumnUpdatedAt,
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
		Where(sq.Eq{ColumnCharacterID: characterID}).
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

	columns := make([]string, 0, len(r.users.columns))
	for _, column := range r.users.columns {
		columns = append(columns, fmt.Sprintf("%s.%s", TableUsers, column))
	}

	query, args, err := sq.Select(columns...).
		From(r.users.table).
		InnerJoin(
			fmt.Sprintf(
				"%s on %s.%s = %s.%s",
				TableCharacters,
				TableCharacters,
				CharacterID,
				TableUsers,
				ColumnCharacterID,
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
		ColumnCharacterID:     user.CharacterID,
		UserAccessToken:       user.AccessToken,
		UserRefreshToken:      user.RefreshToken,
		UserExpires:           user.Expires,
		UserOwnerHash:         user.OwnerHash,
		UserScopes:            user.Scopes,
		UserDisabled:          user.Disabled,
		UserDisabledReason:    user.DisabledReason,
		UserDisabledTimestamp: user.DisabledTimestamp,
		UserLastLogin:         user.LastLogin,
		UserLastProcessed:     user.LastProcessed,
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
		ColumnCharacterID:     user.CharacterID,
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

func (r *userRepository) UsersSortedByProcessedAtLimit(ctx context.Context, limit uint64) ([]*skillz.User, error) {

	query, args, err := sq.Select(r.users.columns...).From(r.users.table).OrderBy(fmt.Sprintf("%s %s", UserLastProcessed, "DESC")).Limit(limit).ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, userRepositoryIdentifier, "UsersSortedByProcessedAtLimit", "failed to generate sql")
	}

	var users = make([]*skillz.User, 0)
	err = r.db.SelectContext(ctx, &users, query, args...)
	return users, err

}

var skillMetaInnerJoin = fmt.Sprintf("%s csm on csm.character_id = users.character_id", TableCharacterSkillMeta)

func (r *userRepository) NewUsersBySP(ctx context.Context) ([]*skillz.User, error) {

	columns := make([]string, 0, len(r.users.columns))
	for _, c := range r.users.columns {
		columns = append(columns, fmt.Sprintf("users.%s", c))
	}

	query, args, err := sq.Select(columns...).
		From(r.users.table).
		InnerJoin(skillMetaInnerJoin).
		Where(fmt.Sprintf("users.%s >= DATE(NOW() - INTERVAL 7 DAY)", ColumnCreatedAt)).
		OrderBy(fmt.Sprintf("users.%s DESC", ColumnCreatedAt)).
		Limit(50).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, userRepositoryIdentifier, "CreatedBySPInLastSevenDays", "failed to generate sql")
	}

	var users = make([]*skillz.User, 0)
	err = r.db.SelectContext(ctx, &users, query, args...)
	return users, err

}

func (r *userRepository) UserSettings(ctx context.Context, id uuid.UUID) (*skillz.UserSettings, error) {

	query, args, err := sq.Select(r.settings.columns...).
		From(r.settings.table).
		Where(sq.Eq{SettingsUserID: id}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, userRepositoryIdentifier, "UserSettings", "failed to generate sql")
	}

	var settings = new(skillz.UserSettings)
	err = r.db.GetContext(ctx, settings, query, args...)
	return settings, errors.Wrapf(err, prefixFormat, userRepositoryIdentifier, "UserSettings")

}

func (r *userRepository) CreateUserSettings(ctx context.Context, settings *skillz.UserSettings) error {

	now := time.Now()
	settings.CreatedAt = now
	settings.UpdatedAt = now

	query, args, err := sq.Insert(r.settings.table).SetMap(map[string]interface{}{
		SettingsUserID:        settings.UserID,
		SettingsHideClones:    settings.HideClones,
		SettingsHideQueue:     settings.HideQueue,
		SettingsHideStandings: settings.HideStandings,
		ColumnCreatedAt:       settings.CreatedAt,
		ColumnUpdatedAt:       settings.UpdatedAt,
	}).
		Suffix(OnDuplicateKeyStmt(SettingsHideClones, SettingsHideQueue, SettingsHideStandings, ColumnUpdatedAt)).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, userRepositoryIdentifier, "CreateUserSettings", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, userRepositoryIdentifier, "CreateUserSettings")

}
