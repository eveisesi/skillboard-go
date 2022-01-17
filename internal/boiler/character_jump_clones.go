// Code generated by SQLBoiler 4.8.3 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package boiler

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// CharacterJumpClone is an object representing the database table.
type CharacterJumpClone struct {
	CharacterID  uint64     `boil:"character_id" json:"character_id" toml:"character_id" yaml:"character_id"`
	JumpCloneID  uint       `boil:"jump_clone_id" json:"jump_clone_id" toml:"jump_clone_id" yaml:"jump_clone_id"`
	LocationID   uint64     `boil:"location_id" json:"location_id" toml:"location_id" yaml:"location_id"`
	LocationType string     `boil:"location_type" json:"location_type" toml:"location_type" yaml:"location_type"`
	ImplantIds   types.JSON `boil:"implant_ids" json:"implant_ids" toml:"implant_ids" yaml:"implant_ids"`
	CreatedAt    time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *characterJumpCloneR `boil:"r" json:"r" toml:"r" yaml:"r"`
	L characterJumpCloneL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CharacterJumpCloneColumns = struct {
	CharacterID  string
	JumpCloneID  string
	LocationID   string
	LocationType string
	ImplantIds   string
	CreatedAt    string
}{
	CharacterID:  "character_id",
	JumpCloneID:  "jump_clone_id",
	LocationID:   "location_id",
	LocationType: "location_type",
	ImplantIds:   "implant_ids",
	CreatedAt:    "created_at",
}

var CharacterJumpCloneTableColumns = struct {
	CharacterID  string
	JumpCloneID  string
	LocationID   string
	LocationType string
	ImplantIds   string
	CreatedAt    string
}{
	CharacterID:  "character_jump_clones.character_id",
	JumpCloneID:  "character_jump_clones.jump_clone_id",
	LocationID:   "character_jump_clones.location_id",
	LocationType: "character_jump_clones.location_type",
	ImplantIds:   "character_jump_clones.implant_ids",
	CreatedAt:    "character_jump_clones.created_at",
}

// Generated where

type whereHelpertypes_JSON struct{ field string }

func (w whereHelpertypes_JSON) EQ(x types.JSON) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertypes_JSON) NEQ(x types.JSON) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertypes_JSON) LT(x types.JSON) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_JSON) LTE(x types.JSON) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_JSON) GT(x types.JSON) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_JSON) GTE(x types.JSON) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var CharacterJumpCloneWhere = struct {
	CharacterID  whereHelperuint64
	JumpCloneID  whereHelperuint
	LocationID   whereHelperuint64
	LocationType whereHelperstring
	ImplantIds   whereHelpertypes_JSON
	CreatedAt    whereHelpertime_Time
}{
	CharacterID:  whereHelperuint64{field: "`character_jump_clones`.`character_id`"},
	JumpCloneID:  whereHelperuint{field: "`character_jump_clones`.`jump_clone_id`"},
	LocationID:   whereHelperuint64{field: "`character_jump_clones`.`location_id`"},
	LocationType: whereHelperstring{field: "`character_jump_clones`.`location_type`"},
	ImplantIds:   whereHelpertypes_JSON{field: "`character_jump_clones`.`implant_ids`"},
	CreatedAt:    whereHelpertime_Time{field: "`character_jump_clones`.`created_at`"},
}

// CharacterJumpCloneRels is where relationship names are stored.
var CharacterJumpCloneRels = struct {
	Character string
}{
	Character: "Character",
}

// characterJumpCloneR is where relationships are stored.
type characterJumpCloneR struct {
	Character *User `boil:"Character" json:"Character" toml:"Character" yaml:"Character"`
}

// NewStruct creates a new relationship struct
func (*characterJumpCloneR) NewStruct() *characterJumpCloneR {
	return &characterJumpCloneR{}
}

// characterJumpCloneL is where Load methods for each relationship are stored.
type characterJumpCloneL struct{}

var (
	characterJumpCloneAllColumns            = []string{"character_id", "jump_clone_id", "location_id", "location_type", "implant_ids", "created_at"}
	characterJumpCloneColumnsWithoutDefault = []string{"character_id", "jump_clone_id", "location_id", "location_type", "implant_ids", "created_at"}
	characterJumpCloneColumnsWithDefault    = []string{}
	characterJumpClonePrimaryKeyColumns     = []string{"character_id", "jump_clone_id"}
)

type (
	// CharacterJumpCloneSlice is an alias for a slice of pointers to CharacterJumpClone.
	// This should almost always be used instead of []CharacterJumpClone.
	CharacterJumpCloneSlice []*CharacterJumpClone
	// CharacterJumpCloneHook is the signature for custom CharacterJumpClone hook methods
	CharacterJumpCloneHook func(context.Context, boil.ContextExecutor, *CharacterJumpClone) error

	characterJumpCloneQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	characterJumpCloneType                 = reflect.TypeOf(&CharacterJumpClone{})
	characterJumpCloneMapping              = queries.MakeStructMapping(characterJumpCloneType)
	characterJumpClonePrimaryKeyMapping, _ = queries.BindMapping(characterJumpCloneType, characterJumpCloneMapping, characterJumpClonePrimaryKeyColumns)
	characterJumpCloneInsertCacheMut       sync.RWMutex
	characterJumpCloneInsertCache          = make(map[string]insertCache)
	characterJumpCloneUpdateCacheMut       sync.RWMutex
	characterJumpCloneUpdateCache          = make(map[string]updateCache)
	characterJumpCloneUpsertCacheMut       sync.RWMutex
	characterJumpCloneUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var characterJumpCloneBeforeInsertHooks []CharacterJumpCloneHook
var characterJumpCloneBeforeUpdateHooks []CharacterJumpCloneHook
var characterJumpCloneBeforeDeleteHooks []CharacterJumpCloneHook
var characterJumpCloneBeforeUpsertHooks []CharacterJumpCloneHook

var characterJumpCloneAfterInsertHooks []CharacterJumpCloneHook
var characterJumpCloneAfterSelectHooks []CharacterJumpCloneHook
var characterJumpCloneAfterUpdateHooks []CharacterJumpCloneHook
var characterJumpCloneAfterDeleteHooks []CharacterJumpCloneHook
var characterJumpCloneAfterUpsertHooks []CharacterJumpCloneHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *CharacterJumpClone) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range characterJumpCloneBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *CharacterJumpClone) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range characterJumpCloneBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *CharacterJumpClone) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range characterJumpCloneBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *CharacterJumpClone) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range characterJumpCloneBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *CharacterJumpClone) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range characterJumpCloneAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *CharacterJumpClone) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range characterJumpCloneAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *CharacterJumpClone) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range characterJumpCloneAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *CharacterJumpClone) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range characterJumpCloneAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *CharacterJumpClone) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range characterJumpCloneAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddCharacterJumpCloneHook registers your hook function for all future operations.
func AddCharacterJumpCloneHook(hookPoint boil.HookPoint, characterJumpCloneHook CharacterJumpCloneHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		characterJumpCloneBeforeInsertHooks = append(characterJumpCloneBeforeInsertHooks, characterJumpCloneHook)
	case boil.BeforeUpdateHook:
		characterJumpCloneBeforeUpdateHooks = append(characterJumpCloneBeforeUpdateHooks, characterJumpCloneHook)
	case boil.BeforeDeleteHook:
		characterJumpCloneBeforeDeleteHooks = append(characterJumpCloneBeforeDeleteHooks, characterJumpCloneHook)
	case boil.BeforeUpsertHook:
		characterJumpCloneBeforeUpsertHooks = append(characterJumpCloneBeforeUpsertHooks, characterJumpCloneHook)
	case boil.AfterInsertHook:
		characterJumpCloneAfterInsertHooks = append(characterJumpCloneAfterInsertHooks, characterJumpCloneHook)
	case boil.AfterSelectHook:
		characterJumpCloneAfterSelectHooks = append(characterJumpCloneAfterSelectHooks, characterJumpCloneHook)
	case boil.AfterUpdateHook:
		characterJumpCloneAfterUpdateHooks = append(characterJumpCloneAfterUpdateHooks, characterJumpCloneHook)
	case boil.AfterDeleteHook:
		characterJumpCloneAfterDeleteHooks = append(characterJumpCloneAfterDeleteHooks, characterJumpCloneHook)
	case boil.AfterUpsertHook:
		characterJumpCloneAfterUpsertHooks = append(characterJumpCloneAfterUpsertHooks, characterJumpCloneHook)
	}
}

// OneG returns a single characterJumpClone record from the query using the global executor.
func (q characterJumpCloneQuery) OneG(ctx context.Context) (*CharacterJumpClone, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single characterJumpClone record from the query.
func (q characterJumpCloneQuery) One(ctx context.Context, exec boil.ContextExecutor) (*CharacterJumpClone, error) {
	o := &CharacterJumpClone{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "boiler: failed to execute a one query for character_jump_clones")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all CharacterJumpClone records from the query using the global executor.
func (q characterJumpCloneQuery) AllG(ctx context.Context) (CharacterJumpCloneSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all CharacterJumpClone records from the query.
func (q characterJumpCloneQuery) All(ctx context.Context, exec boil.ContextExecutor) (CharacterJumpCloneSlice, error) {
	var o []*CharacterJumpClone

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "boiler: failed to assign all query results to CharacterJumpClone slice")
	}

	if len(characterJumpCloneAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all CharacterJumpClone records in the query, and panics on error.
func (q characterJumpCloneQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all CharacterJumpClone records in the query.
func (q characterJumpCloneQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: failed to count character_jump_clones rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q characterJumpCloneQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q characterJumpCloneQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "boiler: failed to check if character_jump_clones exists")
	}

	return count > 0, nil
}

// Character pointed to by the foreign key.
func (o *CharacterJumpClone) Character(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`character_id` = ?", o.CharacterID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "`users`")

	return query
}

// LoadCharacter allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (characterJumpCloneL) LoadCharacter(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCharacterJumpClone interface{}, mods queries.Applicator) error {
	var slice []*CharacterJumpClone
	var object *CharacterJumpClone

	if singular {
		object = maybeCharacterJumpClone.(*CharacterJumpClone)
	} else {
		slice = *maybeCharacterJumpClone.(*[]*CharacterJumpClone)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &characterJumpCloneR{}
		}
		args = append(args, object.CharacterID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &characterJumpCloneR{}
			}

			for _, a := range args {
				if a == obj.CharacterID {
					continue Outer
				}
			}

			args = append(args, obj.CharacterID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`users`),
		qm.WhereIn(`users.character_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for users")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for users")
	}

	if len(characterJumpCloneAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Character = foreign
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CharacterID == foreign.CharacterID {
				local.R.Character = foreign
				break
			}
		}
	}

	return nil
}

// SetCharacterG of the characterJumpClone to the related item.
// Sets o.R.Character to related.
// Adds o to related.R.CharacterCharacterJumpClones.
// Uses the global database handle.
func (o *CharacterJumpClone) SetCharacterG(ctx context.Context, insert bool, related *User) error {
	return o.SetCharacter(ctx, boil.GetContextDB(), insert, related)
}

// SetCharacter of the characterJumpClone to the related item.
// Sets o.R.Character to related.
// Adds o to related.R.CharacterCharacterJumpClones.
func (o *CharacterJumpClone) SetCharacter(ctx context.Context, exec boil.ContextExecutor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `character_jump_clones` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"character_id"}),
		strmangle.WhereClause("`", "`", 0, characterJumpClonePrimaryKeyColumns),
	)
	values := []interface{}{related.CharacterID, o.CharacterID, o.JumpCloneID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CharacterID = related.CharacterID
	if o.R == nil {
		o.R = &characterJumpCloneR{
			Character: related,
		}
	} else {
		o.R.Character = related
	}

	if related.R == nil {
		related.R = &userR{
			CharacterCharacterJumpClones: CharacterJumpCloneSlice{o},
		}
	} else {
		related.R.CharacterCharacterJumpClones = append(related.R.CharacterCharacterJumpClones, o)
	}

	return nil
}

// CharacterJumpClones retrieves all the records using an executor.
func CharacterJumpClones(mods ...qm.QueryMod) characterJumpCloneQuery {
	mods = append(mods, qm.From("`character_jump_clones`"))
	return characterJumpCloneQuery{NewQuery(mods...)}
}

// FindCharacterJumpCloneG retrieves a single record by ID.
func FindCharacterJumpCloneG(ctx context.Context, characterID uint64, jumpCloneID uint, selectCols ...string) (*CharacterJumpClone, error) {
	return FindCharacterJumpClone(ctx, boil.GetContextDB(), characterID, jumpCloneID, selectCols...)
}

// FindCharacterJumpClone retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCharacterJumpClone(ctx context.Context, exec boil.ContextExecutor, characterID uint64, jumpCloneID uint, selectCols ...string) (*CharacterJumpClone, error) {
	characterJumpCloneObj := &CharacterJumpClone{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `character_jump_clones` where `character_id`=? AND `jump_clone_id`=?", sel,
	)

	q := queries.Raw(query, characterID, jumpCloneID)

	err := q.Bind(ctx, exec, characterJumpCloneObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "boiler: unable to select from character_jump_clones")
	}

	if err = characterJumpCloneObj.doAfterSelectHooks(ctx, exec); err != nil {
		return characterJumpCloneObj, err
	}

	return characterJumpCloneObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *CharacterJumpClone) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *CharacterJumpClone) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("boiler: no character_jump_clones provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(characterJumpCloneColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	characterJumpCloneInsertCacheMut.RLock()
	cache, cached := characterJumpCloneInsertCache[key]
	characterJumpCloneInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			characterJumpCloneAllColumns,
			characterJumpCloneColumnsWithDefault,
			characterJumpCloneColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(characterJumpCloneType, characterJumpCloneMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(characterJumpCloneType, characterJumpCloneMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `character_jump_clones` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `character_jump_clones` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `character_jump_clones` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, characterJumpClonePrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "boiler: unable to insert into character_jump_clones")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.CharacterID,
		o.JumpCloneID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "boiler: unable to populate default values for character_jump_clones")
	}

CacheNoHooks:
	if !cached {
		characterJumpCloneInsertCacheMut.Lock()
		characterJumpCloneInsertCache[key] = cache
		characterJumpCloneInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single CharacterJumpClone record using the global executor.
// See Update for more documentation.
func (o *CharacterJumpClone) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the CharacterJumpClone.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *CharacterJumpClone) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	characterJumpCloneUpdateCacheMut.RLock()
	cache, cached := characterJumpCloneUpdateCache[key]
	characterJumpCloneUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			characterJumpCloneAllColumns,
			characterJumpClonePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("boiler: unable to update character_jump_clones, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `character_jump_clones` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, characterJumpClonePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(characterJumpCloneType, characterJumpCloneMapping, append(wl, characterJumpClonePrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to update character_jump_clones row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: failed to get rows affected by update for character_jump_clones")
	}

	if !cached {
		characterJumpCloneUpdateCacheMut.Lock()
		characterJumpCloneUpdateCache[key] = cache
		characterJumpCloneUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q characterJumpCloneQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q characterJumpCloneQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to update all for character_jump_clones")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to retrieve rows affected for character_jump_clones")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o CharacterJumpCloneSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CharacterJumpCloneSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("boiler: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), characterJumpClonePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `character_jump_clones` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, characterJumpClonePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to update all in characterJumpClone slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to retrieve rows affected all in update all characterJumpClone")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *CharacterJumpClone) UpsertG(ctx context.Context, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateColumns, insertColumns)
}

var mySQLCharacterJumpCloneUniqueColumns = []string{}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *CharacterJumpClone) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("boiler: no character_jump_clones provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(characterJumpCloneColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLCharacterJumpCloneUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	characterJumpCloneUpsertCacheMut.RLock()
	cache, cached := characterJumpCloneUpsertCache[key]
	characterJumpCloneUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			characterJumpCloneAllColumns,
			characterJumpCloneColumnsWithDefault,
			characterJumpCloneColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			characterJumpCloneAllColumns,
			characterJumpClonePrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("boiler: unable to upsert character_jump_clones, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`character_jump_clones`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `character_jump_clones` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(characterJumpCloneType, characterJumpCloneMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(characterJumpCloneType, characterJumpCloneMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "boiler: unable to upsert for character_jump_clones")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(characterJumpCloneType, characterJumpCloneMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "boiler: unable to retrieve unique values for character_jump_clones")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "boiler: unable to populate default values for character_jump_clones")
	}

CacheNoHooks:
	if !cached {
		characterJumpCloneUpsertCacheMut.Lock()
		characterJumpCloneUpsertCache[key] = cache
		characterJumpCloneUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single CharacterJumpClone record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *CharacterJumpClone) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single CharacterJumpClone record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *CharacterJumpClone) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("boiler: no CharacterJumpClone provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), characterJumpClonePrimaryKeyMapping)
	sql := "DELETE FROM `character_jump_clones` WHERE `character_id`=? AND `jump_clone_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to delete from character_jump_clones")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: failed to get rows affected by delete for character_jump_clones")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q characterJumpCloneQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q characterJumpCloneQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("boiler: no characterJumpCloneQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to delete all from character_jump_clones")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: failed to get rows affected by deleteall for character_jump_clones")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o CharacterJumpCloneSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CharacterJumpCloneSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(characterJumpCloneBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), characterJumpClonePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `character_jump_clones` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, characterJumpClonePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to delete all from characterJumpClone slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: failed to get rows affected by deleteall for character_jump_clones")
	}

	if len(characterJumpCloneAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *CharacterJumpClone) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("boiler: no CharacterJumpClone provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *CharacterJumpClone) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindCharacterJumpClone(ctx, exec, o.CharacterID, o.JumpCloneID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CharacterJumpCloneSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("boiler: empty CharacterJumpCloneSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CharacterJumpCloneSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CharacterJumpCloneSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), characterJumpClonePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `character_jump_clones`.* FROM `character_jump_clones` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, characterJumpClonePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "boiler: unable to reload all in CharacterJumpCloneSlice")
	}

	*o = slice

	return nil
}

// CharacterJumpCloneExistsG checks if the CharacterJumpClone row exists.
func CharacterJumpCloneExistsG(ctx context.Context, characterID uint64, jumpCloneID uint) (bool, error) {
	return CharacterJumpCloneExists(ctx, boil.GetContextDB(), characterID, jumpCloneID)
}

// CharacterJumpCloneExists checks if the CharacterJumpClone row exists.
func CharacterJumpCloneExists(ctx context.Context, exec boil.ContextExecutor, characterID uint64, jumpCloneID uint) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `character_jump_clones` where `character_id`=? AND `jump_clone_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, characterID, jumpCloneID)
	}
	row := exec.QueryRowContext(ctx, sql, characterID, jumpCloneID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "boiler: unable to check if character_jump_clones exists")
	}

	return exists, nil
}
