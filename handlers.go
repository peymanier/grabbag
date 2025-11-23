package main

import (
	"net/http"
	"text/template"
)

func (s *Server) ListAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := s.Queries.ListAssets(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.gohtml", "templates/assets.gohtml"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = tmpl.Execute(w, AssetsToDTO(assets))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
