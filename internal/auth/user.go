package auth

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type user interface {
	CookieForUserID(ctx context.Context, cookie *http.Cookie, userID uuid.UUID) (*http.Cookie, error)
	UserIDFromCookie(ctx context.Context, cookie *http.Cookie) (uuid.UUID, error)
}

func (s *Service) CookieForUserID(ctx context.Context, cookie *http.Cookie, userID uuid.UUID) (*http.Cookie, error) {

	userIDString := userID.String()

	hash, err := s.hash([]byte(userIDString))
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash user id")
	}

	signature, err := s.userAuth.rsaKey.Sign(rand.Reader, hash, &rsa.PSSOptions{SaltLength: 32, Hash: crypto.SHA256})
	if err != nil {
		return cookie, err
	}

	cookie.Value = fmt.Sprintf("%s.%x", userIDString, signature)
	cookie.Value = userIDString
	cookie.Expires = time.Now().Add(s.userAuth.tokenExpiry)

	fmt.Println(cookie.Value, cookie.Expires.Format("2006-01-02 15:04:05"))

	return cookie, nil

}

func (s *Service) UserIDFromCookie(ctx context.Context, cookie *http.Cookie) (uuid.UUID, error) {

	v := cookie.Value
	vp := strings.Split(v, ".")
	if len(vp) != 2 {
		return uuid.Nil, errors.Errorf("expected split of cookie value to = 2, got %d parts", len(vp))
	}

	userIDStr := vp[0]
	signature := vp[1]

	userIDHash, err := s.hash([]byte(userIDStr))
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to hash user id")
	}

	err = rsa.VerifyPSS(&s.userAuth.rsaKey.PublicKey, crypto.SHA256, []byte(userIDHash), []byte(signature), &rsa.PSSOptions{SaltLength: 32, Hash: crypto.SHA256})
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to verify signature")
	}

	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to parse uuid")
	}

	return userID, nil

}

func (s *Service) hash(in []byte) ([]byte, error) {

	hasher := sha256.New()
	_, err := hasher.Write(in)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil

}
