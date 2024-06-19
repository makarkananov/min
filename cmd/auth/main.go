package main

import (
	"database/sql"
	"errors"
	"flag"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Required for migrations
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq" // Required for postgres
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"min/internal/adapter/handler/grpc/auth"
	"min/internal/adapter/repository/postgres"
	"min/internal/core/service"
	migrations "min/internal/migration"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Parse command line flags
	var configPath string
	var port string
	flag.StringVar(&configPath, "c", "config/auth.yaml", "Path to configuration file")
	flag.StringVar(&port, "p", "50051", "Port to start server on")
	flag.Parse()

	// Initialize and load configuration from file
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Panic("Error loading configuration:", err)
	}

	// Connect to the database and apply migrations
	postgresURL := viper.GetString("postgres_url")
	pgClient, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Panic("Error connecting to the database:", err)
	}
	err = applyMigrations(postgresURL)
	if err != nil {
		log.Panic("Error applying migrations:", err)
	}

	// Create and start the server
	tokenMaxTime := viper.GetInt("token_max_time")
	usersRep := postgres.NewUserRepository(pgClient)
	authService := service.NewAuthService(usersRep, time.Duration(tokenMaxTime)*time.Minute)
	srv := auth.NewServer(authService)
	go func() {
		if err := srv.Start(port); err != nil {
			log.Panicf("Error starting server: %v", err)
		}
	}()
	defer srv.Stop()

	// Set up a signal channel to handle interrupt and termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal
	sig := <-sigCh
	log.Printf("Received signal %v. Shutting down...", sig)
}

// applyMigrations applies all available migrations to the database.
func applyMigrations(dbURL string) error {
	log.Println("Trying to apply migrations...")

	d, err := iofs.New(migrations.FS, "pg/auth")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
