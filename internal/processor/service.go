package processor

import (
	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/user"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger *logrus.Logger
	cache  cache.QueueAPI
	user   user.API

	scopes skillz.ScopeProcessors
}

// type API interface {
// 	Run()
// }

func New(logger *logrus.Logger, cache cache.QueueAPI, user user.API, scopes skillz.ScopeProcessors) *Service {
	return &Service{
		logger: logger,
		cache:  cache,
		user:   user,
		scopes: scopes,
	}
}

func (s *Service) Run() {

}
