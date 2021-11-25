package null

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/volatiletech/null"
)

func MarshalTime(nt null.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if !nt.Valid {
			_, _ = io.WriteString(w, `null`)
			return
		}

		_, _ = io.WriteString(w, strconv.Quote(nt.Time.Format(time.RFC3339)))
	})
}

func UnmarshalTime(i interface{}) (null.Time, error) {

	if i == nil {
		return null.Time{}, nil
	}

	switch i := i.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, i)
		if err != nil {
			return null.Time{}, fmt.Errorf("unable to parse %v as time", i)
		}

		return null.NewTime(t, true), nil
	default:
		return null.Time{}, fmt.Errorf("%v is not a parsable string", i)
	}

}
