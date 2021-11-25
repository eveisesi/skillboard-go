package scalar

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalFloat32(f float32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatFloat(float64(f), 'f', -1, 32))
	})
}

func UnmarshalFloat32(v interface{}) (float32, error) {

	switch v := v.(type) {
	case string:
		f64, e := strconv.ParseFloat(v, 32)
		if e != nil {
			return 0, e
		}
		return float32(f64), nil
	case int:
		return float32(v), nil
	case float32:
		return v, nil
	case json.Number:
		f64, e := strconv.ParseFloat(string(v), 32)
		if e != nil {
			return 0, e
		}
		return float32(f64), nil
	default:
		return 0, fmt.Errorf("%T is not an float", v)
	}

}

func MarshalFloat64(f float64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatFloat(f, 'f', -1, 64))
	})
}

func UnmarshalFloat64(v interface{}) (float64, error) {

	switch v := v.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case json.Number:
		return strconv.ParseFloat(string(v), 64)
	default:
		return 0, fmt.Errorf("%T is not an float", v)
	}

}
