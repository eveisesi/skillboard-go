package view

import (
	"github.com/eveisesi/skillz/internal/graphql"
)

type API interface {
}

type Service struct {
	graphql graphql.API
	// skill   skill.API
}

func New(graphql graphql.API) API {
	return &Service{
		graphql: graphql,
	}
}
