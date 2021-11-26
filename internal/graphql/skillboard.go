package graphql

import (
	"context"
	"encoding/json"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (s *Service) executeQuery(ctx context.Context, params *graphql.RawParams) *graphql.Response {

	ctx = graphql.StartOperationTrace(ctx)

	params.ReadTime = graphql.TraceTiming{
		Start: graphql.Now(),
		End:   graphql.Now(),
	}

	oc, el := s.executor.CreateOperationContext(ctx, params)
	if len(el) > 0 {
		return &graphql.Response{Errors: el}
	}

	var rh graphql.ResponseHandler
	rh, ctx = s.executor.DispatchOperation(ctx, oc)

	return rh(ctx)

}

func (s *Service) Skillboard(ctx context.Context, characterID uint64) (*Skillboard, gqlerror.List) {

	params := &graphql.RawParams{
		Query: skillboard,
		Variables: map[string]interface{}{
			"id": characterID,
		},
	}

	res := s.executeQuery(ctx, params)

	var skillboard = new(Skillboard)
	err := json.Unmarshal(res.Data, skillboard)
	if err != nil {
		return nil, gqlerror.List{gqlerror.Errorf("failed to decode graphql response")}
	}

	return skillboard, nil

}
