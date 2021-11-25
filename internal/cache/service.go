package cache

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Service struct {
	redis *redis.Client
}

const (
	allianceAPI    string = "AllianceAPI"
	authAPI        string = "AuthAPI"
	characterAPI   string = "CharacterAPI"
	cloneAPI       string = "CloneAPI"
	corporationAPI string = "CorporationAPI"
	etagAPI        string = "EtagAPI"
	pageAPI        string = "PageAPI"
	skillAPI       string = "SkillAPI"
	universeAPI    string = "UniverseAPI"
)

const (
	errorFFormat string = "[%s.%s] %s"
)

func New(redis *redis.Client) *Service {
	return &Service{
		redis: redis,
	}
}

func generateKey(args ...string) string {
	return strings.Join(append([]string{"skillz"}, args...), "::")
}

func hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
