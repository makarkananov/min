package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Required for migrations
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq" // Required for postgres
	goredis "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"min/internal/adapter/client/auth"
	handler "min/internal/adapter/handler/http"
	"min/internal/adapter/kafka"
	"min/internal/adapter/repository/postgres"
	"min/internal/adapter/repository/redis"
	"min/internal/core/domain"
	"min/internal/core/service"
	migrations "min/internal/migration"
	"min/pkg/middleware"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Parse command line flags
	var configPath string
	var port string
	flag.StringVar(&configPath, "c", "config/shortener.yaml", "Path to configuration file")
	flag.StringVar(&port, "p", "8080", "Port to start server on")
	flag.Parse()

	// Initialize and load configuration from file
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Panic("Error loading configuration:", err)
	}

	// Add context with cancel function
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize logger
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	// Create a new instance of the Postgres URL repository.
	postgresURL := viper.GetString("postgres_url")
	pgClient, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Panic("Error connecting to the database:", err)
	}
	err = applyMigrations(postgresURL)
	if err != nil {
		log.Panic("Error applying migrations:", err)
	}
	pgRepo := postgres.NewURLRepository(pgClient)

	// Create a new instance of the Redis URL repository. It will be used as a cache.
	opt, err := goredis.ParseURL(viper.GetString("redis_url"))
	if err != nil {
		log.Panic("Error parsing redis url:", err)
	}
	redisClient := goredis.NewClient(opt)
	redisRepo := redis.NewURLRepository(redisClient)

	// Create a new instance of the Auth client.
	authClient, err := auth.NewClient(viper.GetString("auth_server_url"))
	if err != nil {
		log.Panic("Error creating auth client:", err)
	}

	// Create a new instance of the ShortenerService.
	shortenerService := service.NewShortener(pgRepo, redisRepo, viper.GetInt("shorten_length"), authClient)

	// Kafka producer
	kafkaBrokers := viper.GetStringSlice("kafka_brokers")
	kafkaTopics := viper.GetString("kafka_event_topic")
	eventProducer, err := kafka.NewEventProducer(kafkaBrokers, kafkaTopics)
	if err != nil {
		log.Panic("Error creating event producer:", err)
	}

	// Initialize and run the server
	mux := http.NewServeMux()
	shortenerHandler := handler.NewShortenerHandler(shortenerService, eventProducer)
	if err != nil {
		log.Panic("Error creating auth client:", err)
	}
	authHandler := handler.NewAuthHandler(authClient)
	mux.HandleFunc("POST /shorten", middleware.Chain(
		shortenerHandler.Shorten,
		handler.AuthenticationMiddleware(authClient, true),
		handler.AuthorizationMiddleware(domain.USER),
	))
	mux.HandleFunc("DELETE /remove", middleware.Chain(
		shortenerHandler.Remove,
		handler.AuthenticationMiddleware(authClient, true),
		handler.AuthorizationMiddleware(domain.USER),
	))
	mux.HandleFunc("GET /", middleware.Chain(
		shortenerHandler.Redirect,
	))
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /register", middleware.Chain(
		authHandler.Register,
		handler.AuthenticationMiddleware(authClient, true),
		handler.AuthorizationMiddleware(domain.ADMIN),
	))

	rl := middleware.NewRateLimiter(viper.GetInt64("rate_limit"), viper.GetInt64("max_tokens"))
	cl := middleware.NewConcurrencyLimiter(viper.GetInt("concurrency_limit"))

	srv := &http.Server{
		Addr: ":" + port,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		Handler:           middleware.Chain(mux.ServeHTTP, rl.Limit, cl.Limit),
		ReadHeaderTimeout: 5 * time.Second,
	}
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Printf("Server is running on port %s...", port)
		return srv.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Printf("Shut down signal received, shutting down server...")
		return srv.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		log.Println("Exit reason:", err)
	}
}

// applyMigrations applies all available migrations to the database.
func applyMigrations(dbURL string) error {
	log.Println("Trying to apply migrations...")

	d, err := iofs.New(migrations.FS, "pg/shortener")
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
