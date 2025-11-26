package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"apps/api/internal/config"
)

type Service interface {
	Close() error
	GetDB() *pgxpool.Pool
	Health() map[string]string
}

type service struct {
	config *config.DbConfig
	db     *pgxpool.Pool
}

var dbInstance *service

func New(dbConfig *config.DbConfig) Service {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		dbConfig.DbUsername,
		dbConfig.DbPassword,
		dbConfig.DbHost,
		strconv.Itoa(dbConfig.DbPort),
		dbConfig.DbName,
		dbConfig.DbSchema,
	)

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Failed to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	dbInstance = &service{
		config: dbConfig,
		db:     pool,
	}

	return dbInstance
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", s.config.DbName)
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
