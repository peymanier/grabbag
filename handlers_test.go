package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var TestSever *Server

func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "grabbag"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable", "application_name=grabbag")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = RunMigrations(ctx, connStr, "./sql/schema")
	if err != nil {
		log.Fatal(err)
	}

	s := NewServer(conn)
	s.MountHandlers()

	TestSever = s

	exitCode := m.Run()

	if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}

	err = conn.Close(ctx)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitCode)
}

func TestPing(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	assert.NoError(t, err)

	response := executeRequest(req, TestSever)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Body.String(), ".")
}
