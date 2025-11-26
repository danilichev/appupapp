package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"apps/api/internal/config"
	"apps/api/internal/database"
)

type Server struct {
	config *config.Config
	db     database.Service
}

func NewServer() *http.Server {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	NewServer := &Server{
		config: config,
		db:     database.New(config.Db),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.config.App.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
