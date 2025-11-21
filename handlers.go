package main

import (
	"net/http"

	"github.com/peymanier/grabbag/messages"
)

func (s *Server) ListAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := s.Queries.ListAssets(r.Context())
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, messages.ErrUnknown)
	}

	JSONResponse(w, http.StatusOK, AssetsToResponse(assets))
}
