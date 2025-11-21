package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/peymanier/grabbag/database"
)

type Server struct {
	Router  *chi.Mux
	Queries *database.Queries
}

func NewServer(conn database.DBTX) *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	s.Queries = database.New(conn)
	return s
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Heartbeat("/ping"))

	s.Router.Get("/assets", s.ListAssets)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	s := NewServer(conn)
	s.MountHandlers()

	err = http.ListenAndServe(":3333", s.Router)
	if err != nil {
		log.Fatal(err)
	}
}
