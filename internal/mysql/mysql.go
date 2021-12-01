package mysql

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"xorm.io/builder"
)

type tableConf struct {
	table   string
	columns []string
}

func BuildFilters(s sq.SelectBuilder, operators ...*skillz.Operator) sq.SelectBuilder {
	for _, a := range operators {
		if !a.Operation.IsValid() {
			continue
		}

		switch a.Operation {
		case skillz.EqualOp:
			s = s.Where(sq.Eq{a.Column: a.Value})
		case skillz.NotEqualOp:
			s = s.Where(sq.NotEq{a.Column: a.Value})
		case skillz.GreaterThanEqualToOp:
			s = s.Where(sq.GtOrEq{a.Column: a.Value})
		case skillz.GreaterThanOp:
			s = s.Where(sq.Gt{a.Column: a.Value})
		case skillz.LessThanEqualToOp:
			s = s.Where(sq.LtOrEq{a.Column: a.Value})
		case skillz.LessThanOp:
			s = s.Where(sq.Lt{a.Column: a.Value})
		case skillz.InOp:
			s = s.Where(sq.Eq{a.Column: a.Value.(interface{})})
		case skillz.NotInOp:
			s = s.Where(sq.NotEq{a.Column: a.Value.([]interface{})})
		case skillz.LikeOp:
			s = s.Where(sq.Like{a.Column: fmt.Sprintf("%%%v%%", a.Value)})
		case skillz.OrderOp:
			s = s.OrderBy(fmt.Sprintf("%s %s", a.Column, a.Value))
		case skillz.LimitOp:
			s = s.Limit(uint64(a.Value.(int64)))
		case skillz.SkipOp:
			s = s.Offset(uint64(a.Value.(int64)))
		}
	}

	return s

}

func BuildOperators(operators ...*skillz.Operator) *builder.Builder {

	b := builder.MySQL()

	for _, a := range operators {
		if !a.Operation.IsValid() {
			continue
		}

		switch a.Operation {
		case skillz.EqualOp:
			b = b.Where(builder.Eq{a.Column: a.Value})
		case skillz.NotEqualOp:
			b = b.Where(builder.Neq{a.Column: a.Value})
		case skillz.GreaterThanEqualToOp:
			b = b.Where(builder.Gte{a.Column: a.Value})
		case skillz.GreaterThanOp:
			b = b.Where(builder.Gt{a.Column: a.Value})
		case skillz.LessThanEqualToOp:
			b = b.Where(builder.Lte{a.Column: a.Value})
		case skillz.LessThanOp:
			b = b.Where(builder.Lt{a.Column: a.Value})
		case skillz.InOp:
			b = b.Where(builder.In(a.Column, a.Value.([]interface{})...))
		case skillz.NotInOp:
			b = b.Where(builder.NotIn(a.Column, a.Value.([]interface{})...))
		case skillz.LikeOp:
			b = b.Where(builder.Like{a.Column, a.Value.(string)})
		case skillz.OrderOp:
			b = b.OrderBy(fmt.Sprintf("%s %s", a.Column, a.Value))
		case skillz.LimitOp:
			b = b.Limit(a.Value.(int))
		}
	}

	return b

}

func OnDuplicateKeyStmt(columns ...string) string {
	if len(columns) == 0 {
		return ""
	}

	stmts := make([]string, 0, len(columns))
	for _, column := range columns {
		stmts = append(stmts, fmt.Sprintf("%[1]s = VALUES(%[1]s)", column))
	}

	return fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", strings.Join(stmts, ","))
}
