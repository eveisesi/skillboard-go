package cache

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Service struct {
	redis    *redis.Client
	disabled bool
}

const (
	allianceAPI  string = "AllianceAPI"
	authAPI      string = "AuthAPI"
	characterAPI string = "CharacterAPI"
	cloneAPI     string = "CloneAPI"
	// contactAPI     string = "ContactAPI"
	corporationAPI string = "CorporationAPI"
	etagAPI        string = "EtagAPI"
	// pageAPI        string = "PageAPI"
	skillAPI    string = "SkillAPI"
	universeAPI string = "UniverseAPI"
	userAPI     string = "UserAPI"
)

const (
	errorFFormat string = "[%s.%s] %s"
)

func New(redis *redis.Client, disabled bool) *Service {
	return &Service{
		redis:    redis,
		disabled: disabled,
	}
}

func generateKey(args ...string) string {
	return strings.Join(append([]string{"skillz"}, args...), "::")
}

func hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
