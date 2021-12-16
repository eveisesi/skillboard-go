package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func (s *server) handleGetUserContactsByID(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	userID := chi.URLParam(r, "userID")

	user, err := s.user.User(ctx, uuid.FromStringOrNil(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.writeError(ctx, w, http.StatusNotFound, errors.Errorf("failed to find character with uuid of %s", userID))
			return
		}
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error encountered fetching character"))
		return
	}

	contacts, err := s.contacts.Contacts(ctx, user.CharacterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.writeError(ctx, w, http.StatusNotFound, errors.Errorf("failed to fetching skills for character %d", user.CharacterID))
			return
		}
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error encountered fetching character"))
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, contacts)

}
