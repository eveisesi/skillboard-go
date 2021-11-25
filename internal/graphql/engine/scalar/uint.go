package scalar

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalUint(i uint) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatUint(uint64(i), 10))
	})
}

func UnmarshalUint(v interface{}) (uint, error) {
	switch v := v.(type) {
	case string:
		i, e := strconv.ParseUint(v, 10, 32)
		if e != nil {
			return 0, e
		}
		return uint(i), nil
	case uint:
		return v, nil
	case uint64:
		return uint(v), nil
	case json.Number:

		i, e := v.Int64()
		if e != nil {
			return 0, e
		}

		if i < 0 {
			return 0, nil
		}

		return uint(i), nil

	default:
		return 0, fmt.Errorf("%T is not an uint", v)
	}
}

func MarshalUint64(i uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatUint(i, 10))
	})
}

func UnmarshalUint64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case string:

		i, e := strconv.ParseUint(v, 10, 64)
		if e != nil {
			return 0, e
		}
		return uint64(i), nil

	case uint:
		return uint64(v), nil
	case uint64:
		return v, nil
	case json.Number:

		i, e := v.Int64()
		if e != nil {
			return 0, e
		}

		if i < 0 {
			return 0, nil
		}

		return uint64(i), nil

	default:
		return 0, fmt.Errorf("%T is not an uint64", v)
	}
}
