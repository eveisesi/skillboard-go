package cache

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Service struct {
	redis *redis.Client
}

func New(redis *redis.Client) *Service {
	return &Service{
		redis: redis,
	}
}

func generateKey(args ...string) string {
	return strings.Join(append([]string{"skillz"}, args...), "::")
}

func hashUint64(i uint64) string {
	return hash(strconv.FormatUint(i, 10))
}

func hashUint(i uint) string {
	return hash(strconv.FormatUint(uint64(i), 10))
}

func hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
