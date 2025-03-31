package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	Close() error
	GetDB() *pgxpool.Pool
	Health() map[string]string
}

type service struct {
	db *pgxpool.Pool
}

var (
	database   = os.Getenv("DB_NAME")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("failed to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("failed to create connection pool: %v", err)
	}

	dbInstance = &service{
		db: pool,
	}
	return dbInstance
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	s.db.Close()
	return nil
}

func (s *service) GetDB() *pgxpool.Pool {
	return s.db
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"
	return stats
}
