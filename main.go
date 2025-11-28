package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/peymanier/grabbag/database"
	"github.com/peymanier/grabbag/providers"
	"github.com/pressly/goose/v3"
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

func RunMigrations(ctx context.Context, connStr, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %s", err)
	}

	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		return err
	}
	db := stdlib.OpenDB(*config)
	defer db.Close()

	currVersion, err := goose.GetDBVersionContext(ctx, db)
	if err != nil {
		return err
	}

	log.Printf("current migration version: %d", currVersion)

	if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %s", err)
	}

	newVersion, err := goose.GetDBVersionContext(ctx, db)
	if err != nil {
		return err
	}

	if newVersion > currVersion {
		log.Printf("succesfully migrated from version %d to %d", currVersion, newVersion)
	} else {
		log.Printf("database is already up to date at version %d", currVersion)
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	connStr := os.Getenv("POSTGRES_URL")

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	err = RunMigrations(ctx, connStr, "./sql/schema")
	if err != nil {
		log.Fatal(err)
	}

	s := NewServer(conn)
	s.MountHandlers()

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				providers.NobitexUpdate(ctx, s.Queries)
				providers.TGJUUpdateAssets(ctx, s.Queries)
			case <-ctx.Done():
				log.Println("stopping ticker")
				return
			}
		}
	}()

	err = http.ListenAndServe(":3333", s.Router)

	cancel()
	ticker.Stop()

	if err != nil {
		log.Fatal(err)
	}
}
