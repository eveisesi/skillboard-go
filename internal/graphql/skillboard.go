package graphql

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/99designs/gqlgen/graphql"
)

func (s *Service) Skillboard(ctx context.Context, characterID uint64) *Response {

	ctx = graphql.StartOperationTrace(ctx)

	params := &graphql.RawParams{
		Query: skillboard,
		Variables: map[string]interface{}{
			"id": characterID,
		},
		ReadTime: graphql.TraceTiming{
			Start: graphql.Now(),
			End:   graphql.Now(),
		},
	}

	oc, el := s.executor.CreateOperationContext(ctx, params)
	if len(el) > 0 {
		return &Response{Errors: el}
	}

	var rh graphql.ResponseHandler
	rh, ctx = s.executor.DispatchOperation(ctx, oc)

	res := rh(ctx)

	var m = make(map[string]interface{})
	dec := json.NewDecoder(bytes.NewReader(res.Data))
	dec.UseNumber()
	_ = dec.Decode(&m)

	return &Response{Errors: res.Errors, Data: m}

}
