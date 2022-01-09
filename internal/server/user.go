package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func (s *server) handleGetUserByID(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	userID := chi.URLParam(r, "userID")

	user, err := s.user.User(ctx, uuid.FromStringOrNil(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.writeError(ctx, w, http.StatusNotFound, errors.Errorf("failed to find character with uuid of %s", userID))
			return
		}
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error encountered fetch character"))
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, user)

}

func (s *server) handleGetUserSearch(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	query := r.URL.Query().Get("q")
	if query == "" {
		err := errors.New("query param q containing search term is required")
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusBadRequest, err)
		return
	}

	users, err := s.user.SearchUsers(ctx, query)
	if err != nil {
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusBadRequest, err)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, users)

}

func (s *server) handleGetUserCharacterByID(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	userID := chi.URLParam(r, "userID")

	user, err := s.user.User(ctx, uuid.FromStringOrNil(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.writeError(ctx, w, http.StatusNotFound, errors.Errorf("failed to find user with uuid of %s", userID))
			return
		}
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error encountered fetching user"))
		return
	}

	character, err := s.character.Character(ctx, user.CharacterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.writeError(ctx, w, http.StatusNotFound, errors.Errorf("failed to find character with id of %d", user.CharacterID))
			return
		}
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error encountered fetch character"))
		return
	}

	if character.CorporationID > 0 {
		character.Corporation, err = s.corporation.Corporation(ctx, character.CorporationID)
		if err != nil {
			LogEntrySetField(ctx, "error", errors.Wrap(err, "failed to resolve characters corporation id"))
		}
	}

	if character.AllianceID.Valid && character.AllianceID.Uint > 0 {
		character.Alliance, err = s.alliance.Alliance(ctx, character.AllianceID.Uint)
		if err != nil {
			LogEntrySetField(ctx, "error", errors.Wrap(err, "failed to resolve characters corporation id"))
		}
	}

	s.writeResponse(ctx, w, http.StatusOK, character)

}

func (s *server) handleGetNewUsersBySP(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	users, err := s.user.NewUsersBySP(ctx)
	if err != nil {
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusBadRequest, errors.New("unexpected error encountered fetch characters"))
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, users)

}
