package main

import (
	"MyFileExporer/common/env"
	"MyFileExporer/common/logger"
	"MyFileExporer/common/models"
	"MyFileExporer/indexer/internal/batch"
	"MyFileExporer/indexer/internal/crawler"
	"MyFileExporer/indexer/internal/db"
	"MyFileExporer/indexer/internal/queue"
	"MyFileExporer/indexer/internal/repo/database"
	"MyFileExporer/indexer/internal/repo/file"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

// Application is the entry point in the data_provider service
type Application struct {
	Logger      *zap.Logger
	Config      ApplicationConfig
	DBRepo      database.Repo
	FileRepo    file.Repo
	Processor   batch.Processor
	Crawler     crawler.Crawler
	EventsQueue *queue.InMemoryQueue
}

type ApplicationConfig struct {
	PostgresDB      *sql.DB
	PostgresConfig  db.Config
	FileTypesConfig models.FileTypesConfig
}

func main() {
	// Env file setup
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := setup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	crawlerChan := make(chan struct{})

	// Start processor
	go func() {
		err := app.Processor.Run(ctx)
		if err != nil {
			app.Logger.Error(err.Error())
			return
		}
	}()

	// Start crawler
	go func() {
		app.Crawler.Run(ctx)
		crawlerChan <- struct{}{}
	}()

	// Main goroutine waits until the crawler is finish for it to finish
	<-crawlerChan

	for {
		if app.EventsQueue.Length() > 0 {
			time.Sleep(time.Second * 10)
		} else {
			break
		}
	}

	app.Logger.Info("main goroutine finished")
}

func setup() *Application {
	var app Application
	app.Config = ApplicationConfig{}

	processorLogger := logger.InitLogger("./processor.log")
	crawlerLogger := logger.InitLogger("./crawler.log")

	app.Logger = logger.InitLogger("./indexer.log")

	// Database Config
	dbAddress := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.GetString("POSTGRES_USER", "admin"),
		env.GetString("POSTGRES_PASSWORD", "admin_password"),
		env.GetString("POSTGRES_IP", "localhost"),
		env.GetString("POSTGRES_PORT", "5434"),
		env.GetString("POSTGRES_DB", "file_system_database"),
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

	app.Config.PostgresDB = pqDB
	app.Config.PostgresConfig = postgresConfig

	fileTypesConfig, err := models.ParseFileTypesConfig("./../common/file_types_config.json")
	if err != nil {
		panic(err)
	}

	app.Config.FileTypesConfig = fileTypesConfig

	// Database Repository
	dbRepo := database.NewRepo(pqDB, fileTypesConfig)
	app.DBRepo = dbRepo

	// File Repository
	fileRepo := file.NewRepo(fileTypesConfig)
	app.FileRepo = fileRepo

	// Events queue
	eventsQueue := queue.NewQueue()
	app.EventsQueue = eventsQueue

	// Batch Processor
	processor := batch.NewProcessor(app.DBRepo, eventsQueue, processorLogger)
	app.Processor = processor

	// File Crawler
	crawlerConfig := app.loadCrawlerConfig("./config.json")
	fileCrawler := crawler.New(app.FileRepo, eventsQueue, crawlerLogger, crawlerConfig)
	app.Crawler = fileCrawler

	return &app
}

// TODO() config doesn't work
// loadCrawlerConfig reads and parses the JSON config file for indexer component.
func (app *Application) loadCrawlerConfig(filePath string) crawler.Config {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		app.Logger.Panic("Couldn't read indexer config.json", zap.Error(err))
		panic(err)
	}

	var config crawler.Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		app.Logger.Panic("Couldn't parse indexer config.json", zap.Error(err))
		panic(err)
	}

	return config
}
