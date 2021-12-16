package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func (s *server) handleGetUserSkillMetaByID(w http.ResponseWriter, r *http.Request) {

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

	meta, err := s.skills.Meta(ctx, user.CharacterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.writeError(ctx, w, http.StatusNotFound, errors.Errorf("failed to fetching skills for character %d", user.CharacterID))
			return
		}
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error encountered fetching character"))
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, meta)

}
func (s *server) handleGetUserSkillsByID(w http.ResponseWriter, r *http.Request) {

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

	skillsGrouped, err := s.skills.SkillsGrouped(ctx, user.CharacterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.writeError(ctx, w, http.StatusNotFound, errors.Errorf("failed to fetching skills for character %d", user.CharacterID))
			return
		}
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error encountered fetching character"))
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, skillsGrouped)

}

func (s *server) handleGetUserFlyableByID(w http.ResponseWriter, r *http.Request) {

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

	flyable, err := s.skills.Flyable(ctx, user.CharacterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.writeError(ctx, w, http.StatusNotFound, errors.Errorf("failed to fetching skills for character %d", user.CharacterID))
			return
		}
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error encountered fetching character"))
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, flyable)

}

func (s *server) handleGetUserQueueByID(w http.ResponseWriter, r *http.Request) {

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

	queue, err := s.skills.SkillQueue(ctx, user.CharacterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.writeError(ctx, w, http.StatusNotFound, errors.Errorf("failed to fetching skills for character %d", user.CharacterID))
			return
		}
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error encountered fetching character"))
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, queue)

}

func (s *server) handleGetUserAttributesByID(w http.ResponseWriter, r *http.Request) {

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

	attributes, err := s.skills.Attributes(ctx, user.CharacterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.writeError(ctx, w, http.StatusNotFound, errors.Errorf("failed to fetch attributes for character %d", user.CharacterID))
			return
		}
		LogEntrySetField(ctx, "error", err)
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error encountered fetching character"))
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, attributes)

}
