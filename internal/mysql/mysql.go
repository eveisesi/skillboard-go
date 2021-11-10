package mysql

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
)

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
