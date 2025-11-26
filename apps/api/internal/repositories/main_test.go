package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"apps/api/internal/config"
	"apps/api/internal/database"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDbConfig *config.DbConfig
var testDbService database.Service

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "testdb"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := tcpostgres.Run(
		context.Background(),
		"postgres:latest",
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPwd),
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

	testDbConfig = &config.DbConfig{
		DbHost:     dbHost,
		DbPort:     dbPort.Int(),
		DbName:     dbName,
		DbUsername: dbUser,
		DbPassword: dbPwd,
		DbSchema:   "public",
	}

	return dbContainer.Terminate, nil
}

func setupTestDatabase() error {
	// Use standard database/sql DB for migrations
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		testDbConfig.DbUsername,
		testDbConfig.DbPassword,
		testDbConfig.DbHost,
		testDbConfig.DbPort,
		testDbConfig.DbName,
	)

	// Open *sql.DB for migration
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return err
	}

	migrationsPath := "file://../database/migrations"

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		return err
	}

	// Now initialize your app DB service for use in tests
	testDbService = database.New(testDbConfig)
	return nil
}

func cleanupTestDatabase() {
	db := testDbService.GetDB()
	_, _ = db.Exec(
		context.Background(),
		"TRUNCATE TABLE users CASCADE",
	)
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	if err := setupTestDatabase(); err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}

	code := m.Run()

	if teardown != nil {
		if err := teardown(context.Background()); err != nil {
			log.Fatalf("could not teardown postgres container: %v", err)
		}
	}

	log.Printf("Tests completed with exit code %d", code)
}
