package main

import (
	"MyFileExporer/common/env"
	"MyFileExporer/common/logger"
	"MyFileExporer/indexer/internal/batch"
	"MyFileExporer/indexer/internal/db"
	"MyFileExporer/indexer/internal/queue"
	"MyFileExporer/indexer/internal/repo/database"
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
)

// Application is the entry point in the data_provider service
type Application struct {
	Logger *zap.Logger
	Config ApplicationConfig
	DBRepo database.Repo
}

type ApplicationConfig struct {
	PostgresDB     *sql.DB
	PostgresConfig db.Config
}

func main() {
	// Env file setup
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := setup()
	eventsQueue := queue.NewQueue()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processor := batch.NewProcessor(app.DBRepo, eventsQueue)

	// Start processor
	go func() {
		err := processor.Run(ctx)
		if err != nil {
			app.Logger.Error(err.Error())
			return
		}
	}()

}

func setup() *Application {
	var app Application
	app.Config = ApplicationConfig{}

	app.Logger = logger.InitLogger("./indexer.log")

	// Database Config
	dbAddress := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.GetString("POSTGRES_USER", "admin"),
		env.GetString("POSTGRES_PASSWORD", "admin_password"),
		env.GetString("POSTGRES_IP", "localhost"),
		env.GetString("POSTGRES_PORT", "5434"),
		env.GetString("POSTGRES_DB", "companies"),
	)

	postgresConfig := db.Config{
		Addr:         dbAddress,
		MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	// Database
	pqDB, err := db.New(postgresConfig)
	if err != nil {
		panic(err)
	}

	app.Logger.Info("Postgres connection pool established")

	// Repository
	dbRepo := database.NewRepo(pqDB)

	app.Config.PostgresDB = pqDB
	app.Config.PostgresConfig = postgresConfig
	app.DBRepo = dbRepo

	return &app
}
