package scalar

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
)

func MarshalUUID(s uuid.UUID) graphql.Marshaler {
	return graphql.WriterFunc(func(writer io.Writer) {
		_, _ = io.WriteString(writer, strconv.Quote(s.String()))
	})
}

func UnmarshalUUID(i interface{}) (uuid.UUID, error) {
	switch v := i.(type) {
	case string:
		return uuid.FromString(v)
	default:
		return uuid.Nil, fmt.Errorf("%v is not a parsable uuid", v)
	}

}
