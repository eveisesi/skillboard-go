package null

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/volatiletech/null"
)

func MarshalString(ns null.String) graphql.Marshaler {

	if !ns.Valid {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {

		_, _ = io.WriteString(w, strconv.Quote(ns.String))
	})
}

func UnmarshalString(i interface{}) (null.String, error) {
	switch v := i.(type) {
	case string:
		if v == "null" {
			return null.NewString("", false), nil
		}
		return null.NewString(v, true), nil
	default:
		return null.NewString("", false), fmt.Errorf("%v is not a valid string", v)
	}
}
