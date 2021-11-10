package skillz

import (
	"github.com/volatiletech/null"
)

type AuthAttempt struct {
	Status AuthAttemptStatus
	State  string
	Token  null.String
	User   *User
}

type AuthAttemptStatus uint

const (
	CreateAuthStatus AuthAttemptStatus = iota
	PendingAuthStatus
	InvalidAuthStatus
	ExpiredAuthStatus
	CompletedAuthStatus
)
