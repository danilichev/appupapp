package database

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"apps/api/internal/config"
)

var dbConfig *config.DbConfig

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "testdb"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbConfig = &config.DbConfig{
		DbHost:     dbHost,
		DbPort:     dbPort.Int(),
		DbName:     dbName,
		DbUsername: dbUser,
		DbPassword: dbPwd,
		DbSchema:   "public",
	}

	return dbContainer.Terminate, nil
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	code := m.Run()

	if teardown != nil {
		if err := teardown(context.Background()); err != nil {
			log.Fatalf("could not teardown postgres container: %v", err)
		}
	}

	// Exit with the test run's exit code
	log.Printf("Tests completed with exit code %d", code)
}

func TestNew(t *testing.T) {
	srv := New(dbConfig)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New(dbConfig)

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf(
			"expected message to be 'It's healthy', got %s",
			stats["message"],
		)
	}
}

func TestClose(t *testing.T) {
	srv := New(dbConfig)

	if err := srv.Close(); err != nil {
		t.Fatalf("expected Close() to return nil, got %v", err)
	}
}
