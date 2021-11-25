package null

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/volatiletech/null"
)

func UnmarshalFloat32(f interface{}) (null.Float32, error) {
	switch v := f.(type) {
	case float32:
		return null.NewFloat32(v, true), nil
	case float64:
		return null.NewFloat32(float32(v), true), nil
	default:
		return null.Float32{}, fmt.Errorf("%v is not a valid float32", v)
	}
}

func MarshalFloat32(nf null.Float32) graphql.Marshaler {

	if !nf.Valid {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatFloat(float64(nf.Float32), 'f', -1, 64))
	})

}

func UnmarshalFloat64(f interface{}) (null.Float64, error) {
	switch v := f.(type) {
	case float32:
		return null.NewFloat64(float64(v), true), nil
	case float64:
		return null.NewFloat64(v, true), nil
	default:
		return null.Float64{}, fmt.Errorf("%v is not a valid float32", v)
	}
}

func MarshalFloat64(nf null.Float64) graphql.Marshaler {

	if !nf.Valid {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatFloat(float64(nf.Float64), 'f', -1, 64))
	})

}
