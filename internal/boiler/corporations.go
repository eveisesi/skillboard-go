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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Corporation is an object representing the database table.
type Corporation struct {
	ID            uint        `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name          string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Ticker        string      `boil:"ticker" json:"ticker" toml:"ticker" yaml:"ticker"`
	CeoID         uint        `boil:"ceo_id" json:"ceo_id" toml:"ceo_id" yaml:"ceo_id"`
	CreatorID     uint        `boil:"creator_id" json:"creator_id" toml:"creator_id" yaml:"creator_id"`
	AllianceID    null.Uint   `boil:"alliance_id" json:"alliance_id,omitempty" toml:"alliance_id" yaml:"alliance_id,omitempty"`
	HomeStationID null.Uint   `boil:"home_station_id" json:"home_station_id,omitempty" toml:"home_station_id" yaml:"home_station_id,omitempty"`
	FactionID     null.Uint   `boil:"faction_id" json:"faction_id,omitempty" toml:"faction_id" yaml:"faction_id,omitempty"`
	MemberCount   uint        `boil:"member_count" json:"member_count" toml:"member_count" yaml:"member_count"`
	Shares        null.Uint64 `boil:"shares" json:"shares,omitempty" toml:"shares" yaml:"shares,omitempty"`
	TaxRate       float32     `boil:"tax_rate" json:"tax_rate" toml:"tax_rate" yaml:"tax_rate"`
	URL           null.String `boil:"url" json:"url,omitempty" toml:"url" yaml:"url,omitempty"`
	WarEligible   int8        `boil:"war_eligible" json:"war_eligible" toml:"war_eligible" yaml:"war_eligible"`
	DateFounded   null.Time   `boil:"date_founded" json:"date_founded,omitempty" toml:"date_founded" yaml:"date_founded,omitempty"`
	CreatedAt     time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt     time.Time   `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *corporationR `boil:"r" json:"r" toml:"r" yaml:"r"`
	L corporationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CorporationColumns = struct {
	ID            string
	Name          string
	Ticker        string
	CeoID         string
	CreatorID     string
	AllianceID    string
	HomeStationID string
	FactionID     string
	MemberCount   string
	Shares        string
	TaxRate       string
	URL           string
	WarEligible   string
	DateFounded   string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "id",
	Name:          "name",
	Ticker:        "ticker",
	CeoID:         "ceo_id",
	CreatorID:     "creator_id",
	AllianceID:    "alliance_id",
	HomeStationID: "home_station_id",
	FactionID:     "faction_id",
	MemberCount:   "member_count",
	Shares:        "shares",
	TaxRate:       "tax_rate",
	URL:           "url",
	WarEligible:   "war_eligible",
	DateFounded:   "date_founded",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

var CorporationTableColumns = struct {
	ID            string
	Name          string
	Ticker        string
	CeoID         string
	CreatorID     string
	AllianceID    string
	HomeStationID string
	FactionID     string
	MemberCount   string
	Shares        string
	TaxRate       string
	URL           string
	WarEligible   string
	DateFounded   string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "corporations.id",
	Name:          "corporations.name",
	Ticker:        "corporations.ticker",
	CeoID:         "corporations.ceo_id",
	CreatorID:     "corporations.creator_id",
	AllianceID:    "corporations.alliance_id",
	HomeStationID: "corporations.home_station_id",
	FactionID:     "corporations.faction_id",
	MemberCount:   "corporations.member_count",
	Shares:        "corporations.shares",
	TaxRate:       "corporations.tax_rate",
	URL:           "corporations.url",
	WarEligible:   "corporations.war_eligible",
	DateFounded:   "corporations.date_founded",
	CreatedAt:     "corporations.created_at",
	UpdatedAt:     "corporations.updated_at",
}

// Generated where

type whereHelpernull_Uint64 struct{ field string }

func (w whereHelpernull_Uint64) EQ(x null.Uint64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Uint64) NEQ(x null.Uint64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Uint64) LT(x null.Uint64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Uint64) LTE(x null.Uint64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Uint64) GT(x null.Uint64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Uint64) GTE(x null.Uint64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_Uint64) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Uint64) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

type whereHelperfloat32 struct{ field string }

func (w whereHelperfloat32) EQ(x float32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperfloat32) NEQ(x float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperfloat32) LT(x float32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperfloat32) LTE(x float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperfloat32) GT(x float32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperfloat32) GTE(x float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelperfloat32) IN(slice []float32) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperfloat32) NIN(slice []float32) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var CorporationWhere = struct {
	ID            whereHelperuint
	Name          whereHelperstring
	Ticker        whereHelperstring
	CeoID         whereHelperuint
	CreatorID     whereHelperuint
	AllianceID    whereHelpernull_Uint
	HomeStationID whereHelpernull_Uint
	FactionID     whereHelpernull_Uint
	MemberCount   whereHelperuint
	Shares        whereHelpernull_Uint64
	TaxRate       whereHelperfloat32
	URL           whereHelpernull_String
	WarEligible   whereHelperint8
	DateFounded   whereHelpernull_Time
	CreatedAt     whereHelpertime_Time
	UpdatedAt     whereHelpertime_Time
}{
	ID:            whereHelperuint{field: "`corporations`.`id`"},
	Name:          whereHelperstring{field: "`corporations`.`name`"},
	Ticker:        whereHelperstring{field: "`corporations`.`ticker`"},
	CeoID:         whereHelperuint{field: "`corporations`.`ceo_id`"},
	CreatorID:     whereHelperuint{field: "`corporations`.`creator_id`"},
	AllianceID:    whereHelpernull_Uint{field: "`corporations`.`alliance_id`"},
	HomeStationID: whereHelpernull_Uint{field: "`corporations`.`home_station_id`"},
	FactionID:     whereHelpernull_Uint{field: "`corporations`.`faction_id`"},
	MemberCount:   whereHelperuint{field: "`corporations`.`member_count`"},
	Shares:        whereHelpernull_Uint64{field: "`corporations`.`shares`"},
	TaxRate:       whereHelperfloat32{field: "`corporations`.`tax_rate`"},
	URL:           whereHelpernull_String{field: "`corporations`.`url`"},
	WarEligible:   whereHelperint8{field: "`corporations`.`war_eligible`"},
	DateFounded:   whereHelpernull_Time{field: "`corporations`.`date_founded`"},
	CreatedAt:     whereHelpertime_Time{field: "`corporations`.`created_at`"},
	UpdatedAt:     whereHelpertime_Time{field: "`corporations`.`updated_at`"},
}

// CorporationRels is where relationship names are stored.
var CorporationRels = struct {
}{}

// corporationR is where relationships are stored.
type corporationR struct {
}

// NewStruct creates a new relationship struct
func (*corporationR) NewStruct() *corporationR {
	return &corporationR{}
}

// corporationL is where Load methods for each relationship are stored.
type corporationL struct{}

var (
	corporationAllColumns            = []string{"id", "name", "ticker", "ceo_id", "creator_id", "alliance_id", "home_station_id", "faction_id", "member_count", "shares", "tax_rate", "url", "war_eligible", "date_founded", "created_at", "updated_at"}
	corporationColumnsWithoutDefault = []string{"id", "name", "ticker", "ceo_id", "creator_id", "alliance_id", "home_station_id", "faction_id", "member_count", "shares", "tax_rate", "url", "date_founded", "created_at", "updated_at"}
	corporationColumnsWithDefault    = []string{"war_eligible"}
	corporationPrimaryKeyColumns     = []string{"id"}
)

type (
	// CorporationSlice is an alias for a slice of pointers to Corporation.
	// This should almost always be used instead of []Corporation.
	CorporationSlice []*Corporation
	// CorporationHook is the signature for custom Corporation hook methods
	CorporationHook func(context.Context, boil.ContextExecutor, *Corporation) error

	corporationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	corporationType                 = reflect.TypeOf(&Corporation{})
	corporationMapping              = queries.MakeStructMapping(corporationType)
	corporationPrimaryKeyMapping, _ = queries.BindMapping(corporationType, corporationMapping, corporationPrimaryKeyColumns)
	corporationInsertCacheMut       sync.RWMutex
	corporationInsertCache          = make(map[string]insertCache)
	corporationUpdateCacheMut       sync.RWMutex
	corporationUpdateCache          = make(map[string]updateCache)
	corporationUpsertCacheMut       sync.RWMutex
	corporationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var corporationBeforeInsertHooks []CorporationHook
var corporationBeforeUpdateHooks []CorporationHook
var corporationBeforeDeleteHooks []CorporationHook
var corporationBeforeUpsertHooks []CorporationHook

var corporationAfterInsertHooks []CorporationHook
var corporationAfterSelectHooks []CorporationHook
var corporationAfterUpdateHooks []CorporationHook
var corporationAfterDeleteHooks []CorporationHook
var corporationAfterUpsertHooks []CorporationHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Corporation) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range corporationBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Corporation) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range corporationBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Corporation) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range corporationBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Corporation) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range corporationBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Corporation) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range corporationAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Corporation) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range corporationAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Corporation) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range corporationAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Corporation) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range corporationAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Corporation) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range corporationAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddCorporationHook registers your hook function for all future operations.
func AddCorporationHook(hookPoint boil.HookPoint, corporationHook CorporationHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		corporationBeforeInsertHooks = append(corporationBeforeInsertHooks, corporationHook)
	case boil.BeforeUpdateHook:
		corporationBeforeUpdateHooks = append(corporationBeforeUpdateHooks, corporationHook)
	case boil.BeforeDeleteHook:
		corporationBeforeDeleteHooks = append(corporationBeforeDeleteHooks, corporationHook)
	case boil.BeforeUpsertHook:
		corporationBeforeUpsertHooks = append(corporationBeforeUpsertHooks, corporationHook)
	case boil.AfterInsertHook:
		corporationAfterInsertHooks = append(corporationAfterInsertHooks, corporationHook)
	case boil.AfterSelectHook:
		corporationAfterSelectHooks = append(corporationAfterSelectHooks, corporationHook)
	case boil.AfterUpdateHook:
		corporationAfterUpdateHooks = append(corporationAfterUpdateHooks, corporationHook)
	case boil.AfterDeleteHook:
		corporationAfterDeleteHooks = append(corporationAfterDeleteHooks, corporationHook)
	case boil.AfterUpsertHook:
		corporationAfterUpsertHooks = append(corporationAfterUpsertHooks, corporationHook)
	}
}

// OneG returns a single corporation record from the query using the global executor.
func (q corporationQuery) OneG(ctx context.Context) (*Corporation, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single corporation record from the query.
func (q corporationQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Corporation, error) {
	o := &Corporation{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "boiler: failed to execute a one query for corporations")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all Corporation records from the query using the global executor.
func (q corporationQuery) AllG(ctx context.Context) (CorporationSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all Corporation records from the query.
func (q corporationQuery) All(ctx context.Context, exec boil.ContextExecutor) (CorporationSlice, error) {
	var o []*Corporation

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "boiler: failed to assign all query results to Corporation slice")
	}

	if len(corporationAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all Corporation records in the query, and panics on error.
func (q corporationQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all Corporation records in the query.
func (q corporationQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: failed to count corporations rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q corporationQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q corporationQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "boiler: failed to check if corporations exists")
	}

	return count > 0, nil
}

// Corporations retrieves all the records using an executor.
func Corporations(mods ...qm.QueryMod) corporationQuery {
	mods = append(mods, qm.From("`corporations`"))
	return corporationQuery{NewQuery(mods...)}
}

// FindCorporationG retrieves a single record by ID.
func FindCorporationG(ctx context.Context, iD uint, selectCols ...string) (*Corporation, error) {
	return FindCorporation(ctx, boil.GetContextDB(), iD, selectCols...)
}

// FindCorporation retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCorporation(ctx context.Context, exec boil.ContextExecutor, iD uint, selectCols ...string) (*Corporation, error) {
	corporationObj := &Corporation{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `corporations` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, corporationObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "boiler: unable to select from corporations")
	}

	if err = corporationObj.doAfterSelectHooks(ctx, exec); err != nil {
		return corporationObj, err
	}

	return corporationObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Corporation) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Corporation) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("boiler: no corporations provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(corporationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	corporationInsertCacheMut.RLock()
	cache, cached := corporationInsertCache[key]
	corporationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			corporationAllColumns,
			corporationColumnsWithDefault,
			corporationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(corporationType, corporationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(corporationType, corporationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `corporations` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `corporations` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `corporations` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, corporationPrimaryKeyColumns))
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
		return errors.Wrap(err, "boiler: unable to insert into corporations")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "boiler: unable to populate default values for corporations")
	}

CacheNoHooks:
	if !cached {
		corporationInsertCacheMut.Lock()
		corporationInsertCache[key] = cache
		corporationInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single Corporation record using the global executor.
// See Update for more documentation.
func (o *Corporation) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the Corporation.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Corporation) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	corporationUpdateCacheMut.RLock()
	cache, cached := corporationUpdateCache[key]
	corporationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			corporationAllColumns,
			corporationPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("boiler: unable to update corporations, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `corporations` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, corporationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(corporationType, corporationMapping, append(wl, corporationPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "boiler: unable to update corporations row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: failed to get rows affected by update for corporations")
	}

	if !cached {
		corporationUpdateCacheMut.Lock()
		corporationUpdateCache[key] = cache
		corporationUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q corporationQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q corporationQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to update all for corporations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to retrieve rows affected for corporations")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o CorporationSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CorporationSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), corporationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `corporations` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, corporationPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to update all in corporation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to retrieve rows affected all in update all corporation")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Corporation) UpsertG(ctx context.Context, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateColumns, insertColumns)
}

var mySQLCorporationUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Corporation) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("boiler: no corporations provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(corporationColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLCorporationUniqueColumns, o)

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

	corporationUpsertCacheMut.RLock()
	cache, cached := corporationUpsertCache[key]
	corporationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			corporationAllColumns,
			corporationColumnsWithDefault,
			corporationColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			corporationAllColumns,
			corporationPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("boiler: unable to upsert corporations, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`corporations`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `corporations` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(corporationType, corporationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(corporationType, corporationMapping, ret)
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
		return errors.Wrap(err, "boiler: unable to upsert for corporations")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(corporationType, corporationMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "boiler: unable to retrieve unique values for corporations")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "boiler: unable to populate default values for corporations")
	}

CacheNoHooks:
	if !cached {
		corporationUpsertCacheMut.Lock()
		corporationUpsertCache[key] = cache
		corporationUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single Corporation record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Corporation) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single Corporation record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Corporation) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("boiler: no Corporation provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), corporationPrimaryKeyMapping)
	sql := "DELETE FROM `corporations` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to delete from corporations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: failed to get rows affected by delete for corporations")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q corporationQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q corporationQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("boiler: no corporationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to delete all from corporations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: failed to get rows affected by deleteall for corporations")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o CorporationSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CorporationSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(corporationBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), corporationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `corporations` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, corporationPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "boiler: unable to delete all from corporation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "boiler: failed to get rows affected by deleteall for corporations")
	}

	if len(corporationAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Corporation) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("boiler: no Corporation provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Corporation) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindCorporation(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CorporationSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("boiler: empty CorporationSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CorporationSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CorporationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), corporationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `corporations`.* FROM `corporations` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, corporationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "boiler: unable to reload all in CorporationSlice")
	}

	*o = slice

	return nil
}

// CorporationExistsG checks if the Corporation row exists.
func CorporationExistsG(ctx context.Context, iD uint) (bool, error) {
	return CorporationExists(ctx, boil.GetContextDB(), iD)
}

// CorporationExists checks if the Corporation row exists.
func CorporationExists(ctx context.Context, exec boil.ContextExecutor, iD uint) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `corporations` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "boiler: unable to check if corporations exists")
	}

	return exists, nil
}
