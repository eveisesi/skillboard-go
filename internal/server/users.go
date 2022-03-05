package server

import (
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

func (s *Server) handleGetRecent(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()
	var nr = newrelic.FromContext(ctx)

	highlighted, recent, err := s.users.Recent(ctx)
	if err != nil {
		nr.NoticeError(err)
		s.logger.WithError(err).Error("failed to fetch recent users")
		s.writeError(ctx, w, http.StatusInternalServerError, errors.New("failed to fetch recent users"))
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, struct {
		Recent      []*skillz.User `json:"recent"`
		Highlighted []*skillz.User `json:"highlighted"`
	}{recent, highlighted})

}
