package auth

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type user interface {
	UserCookie(ctx context.Context, userID uuid.UUID) (*http.Cookie, error)
	UserIDFromCookie(ctx context.Context, cookie *http.Cookie) (uuid.UUID, error)
	LogoutCookie(ctx context.Context) (*http.Cookie, error)
}

func (s *Service) UserCookie(ctx context.Context, userID uuid.UUID) (*http.Cookie, error) {

	hash, err := s.hash([]byte(userID.String()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash user id")
	}

	signature, err := s.userAuth.rsaKey.Sign(rand.Reader, hash, &rsa.PSSOptions{Hash: crypto.SHA256})
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:   internal.CookieID,
		Domain: fmt.Sprintf(".%s", s.userAuth.cookieDomain),
		MaxAge: int(s.userAuth.cookieExpiry.Seconds()),
		Path:   "/",
		Value:  fmt.Sprintf("%s.%s", userID, hex.EncodeToString(signature)),
		Secure: s.env == skillz.Production,
	}, nil

}

func (s *Service) LogoutCookie(ctx context.Context) (*http.Cookie, error) {
	return &http.Cookie{
		Name:   internal.CookieID,
		Domain: s.userAuth.cookieDomain,
		MaxAge: -1,
		Path:   "/",
	}, nil
}

func (s *Service) UserIDFromCookie(ctx context.Context, cookie *http.Cookie) (uuid.UUID, error) {

	v := cookie.Value

	// split value it parts seperated by a period
	vp := strings.Split(v, ".")
	if len(vp) != 2 {
		return uuid.Nil, errors.Errorf("expected split of cookie value to = 2, got %d parts", len(vp))
	}

	// The first part should be a valid uuid.
	userID, err := uuid.FromString(vp[0])
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to validate user uuid in cookie")
	}

	// Second part is the integrity hash
	signatureStr := vp[1]

	// Decode the signature from a string to a slice of bytes representing the hex string
	signature, err := hex.DecodeString(signatureStr)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to decode signature string to hex")
	}

	// Hash the User ID
	userIDHash, err := s.hash([]byte(userID.String()))
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to hash user id")
	}

	// Ensure the integrity of the string. This prevents somebody from modifying the value of the cookie and impresonating the user
	err = rsa.VerifyPSS(&s.userAuth.rsaKey.PublicKey, crypto.SHA256, []byte(userIDHash), []byte(signature), &rsa.PSSOptions{Hash: crypto.SHA256})
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to verify signature")
	}

	return userID, nil

}

func (s *Service) hash(in []byte) ([]byte, error) {

	hasher := crypto.SHA256.New()
	_, err := hasher.Write(in)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil

}
