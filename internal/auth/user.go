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

	signature, err := s.userAuth.rsaKey.Sign(rand.Reader, hash, &rsa.PSSOptions{Hash: crypto.SHA256})
	if err != nil {
		return cookie, err
	}

	cookie.Value = fmt.Sprintf("%s.%s", userIDString, hex.EncodeToString(signature))
	cookie.Expires = time.Now().Add(s.userAuth.tokenExpiry)

	return cookie, nil

}

func (s *Service) UserIDFromCookie(ctx context.Context, cookie *http.Cookie) (uuid.UUID, error) {

	v := cookie.Value
	vp := strings.Split(v, ".")
	if len(vp) != 2 {
		return uuid.Nil, errors.Errorf("expected split of cookie value to = 2, got %d parts", len(vp))
	}

	userIDStr := vp[0]
	signatureStr := vp[1]

	signature, err := hex.DecodeString(signatureStr)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to decode signature string to hex")
	}

	userIDHash, err := s.hash([]byte(userIDStr))
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to hash user id")
	}

	err = rsa.VerifyPSS(&s.userAuth.rsaKey.PublicKey, crypto.SHA256, []byte(userIDHash), []byte(signature), &rsa.PSSOptions{Hash: crypto.SHA256})
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

	hasher := crypto.SHA256.New()
	_, err := hasher.Write(in)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil

}
