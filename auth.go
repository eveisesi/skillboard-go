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

var StatusMap = map[AuthAttemptStatus]string{
	CreateAuthStatus:    "create",
	PendingAuthStatus:   "pending",
	InvalidAuthStatus:   "invalid",
	ExpiredAuthStatus:   "expired",
	CompletedAuthStatus: "completed",
}
