package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/peymanier/grabbag/messages"
)

func (s *Server) ListAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := s.Queries.ListAssets(r.Context())
	if err != nil {
		log.Println(err.Error())
		http.Error(w, messages.ErrUnknown.String(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.gohtml", "templates/assets.gohtml"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = tmpl.Execute(w, AssetsToDTO(assets))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, messages.ErrUnknown.String(), http.StatusInternalServerError)
	}
}
