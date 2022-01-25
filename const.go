package skillz

import "strings"

type Environment uint

const (
	Production Environment = iota
	Development
)

func EnvironmentFromString(env string) Environment {
	switch strings.ToLower(env) {
	case "prod", "production":
		return Production
	default:
		return Development
	}
}

func (e Environment) String() string {
	switch e {
	case Production:
		return "production"
	default:
		return "development"
	}
}
