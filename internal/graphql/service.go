package graphql

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/clone"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/eveisesi/skillz/internal/graphql/engine"
	"github.com/eveisesi/skillz/internal/graphql/engine/dataloaders"
	"github.com/eveisesi/skillz/internal/graphql/engine/resolvers"
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/eveisesi/skillz/internal/universe"
	"github.com/eveisesi/skillz/internal/user"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type API interface {
	Skillboard(ctx context.Context, characterID uint64) (*Skillboard, gqlerror.List)
	ExecutableSchema() graphql.ExecutableSchema
}

type Service struct {
	es       graphql.ExecutableSchema
	executor graphql.GraphExecutor
}

func New(
	alliance alliance.API,
	auth auth.API,
	character character.API,
	clone clone.API,
	corporation corporation.API,
	skill skill.API,
	universe universe.API,
	user user.API,
) *Service {

	dl := dataloaders.New(
		time.Millisecond*100, 100, character, corporation, alliance, clone, universe,
	)

	es := engine.NewExecutableSchema(engine.Config{
		Resolvers: resolvers.New(alliance, auth, character, clone, corporation, dl, skill, user),
		Directives: engine.DirectiveRoot{
			IsAuthed: IsAuthed,
		},
	})

	return &Service{es, executor.New(es)}
}

func (s *Service) ExecutableSchema() graphql.ExecutableSchema {
	return s.es
}
