package skill

import (
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/etag"
)

type Service struct {
	cache cache.SkillAPI
	etag  etag.API
	esi   esi.API
}
