package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"min/internal/adapter/kafka"
	"min/internal/migration"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ClickHouse/clickhouse-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	ch "min/internal/adapter/repository/clickhouse"
	"min/internal/core/service"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config/statistics.yaml", "Path to configuration file")
	flag.Parse()

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Panic("Error loading configuration:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	// ClickHouse repository
	clickhouseURL := viper.GetString("clickhouse_url")
	db, err := sql.Open("clickhouse", clickhouseURL)
	if err != nil {
		log.Panicf("Error connecting to ClickHouse: %v", err)
	}
	err = applyMigrations(clickhouseURL)
	if err != nil {
		log.Panic("Error applying migrations:", err)
	}
	chRepo := ch.NewEventRepository(db)

	// Statistics service
	statsService := service.NewStatisticsService(chRepo)

	// Kafka consumer
	kafkaBrokers := viper.GetStringSlice("kafka_brokers")
	consumerGroupID := viper.GetString("kafka_consumer_group_id")
	kafkaTopics := viper.GetStringSlice("kafka_topics")
	kafkaConsumer, err := kafka.NewKafkaConsumer(kafkaBrokers, consumerGroupID, statsService)
	if err != nil {
		log.Panic("Error creating Kafka consumer:", err)
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Println("Starting Kafka consumer...")
		return kafkaConsumer.Start(gCtx, kafkaTopics)
	})

	if err := g.Wait(); err != nil {
		log.Println("Exit reason:", err)
	}
}

// applyMigrations applies all available migrations to the database.
func applyMigrations(dbURL string) error {
	log.Println("Trying to apply migrations...")

	d, err := iofs.New(migrations.FS, "clickhouse/event")
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

	log.Infof("Migrations applied successfully")
	return nil
}
