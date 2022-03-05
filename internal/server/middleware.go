package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (s *Server) monitoring(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		txn := s.newrelic.StartTransaction(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		txn.SetWebRequestHTTP(r)
		rw := txn.SetWebResponse(w)
		defer txn.End()

		r = newrelic.RequestWithTransactionContext(r, txn)

		next.ServeHTTP(rw, r)

		rctx := chi.RouteContext(r.Context())
		name := rctx.RoutePattern()

		// ignore invalid routes
		if name == "/*" {
			txn.Ignore()
			return
		}

		txn.SetName(r.Method + " " + name)

	})
}

func (s *Server) requestLogger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&structuredLogger{s.logger})
}
