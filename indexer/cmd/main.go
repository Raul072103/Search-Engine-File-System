package main

import (
	"MyFileExporer/common/env"
	"MyFileExporer/common/logger"
	"MyFileExporer/indexer/internal/batch"
	"MyFileExporer/indexer/internal/crawler"
	"MyFileExporer/indexer/internal/db"
	"MyFileExporer/indexer/internal/queue"
	"MyFileExporer/indexer/internal/repo/database"
	"MyFileExporer/indexer/internal/repo/file"
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
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

	app.Config.PostgresDB = pqDB
	app.Config.PostgresConfig = postgresConfig

	// Database Repository
	dbRepo := database.NewRepo(pqDB)
	app.DBRepo = dbRepo

	// File Repository
	fileRepo := file.NewRepo()
	app.FileRepo = fileRepo

	// Events queue
	eventsQueue := queue.NewQueue()
	app.EventsQueue = eventsQueue

	// Batch Processor
	processor := batch.NewProcessor(app.DBRepo, eventsQueue, app.Logger)
	app.Processor = processor

	// File Crawler
	crawlerConfig := crawler.Config{
		IgnorePatterns: nil,
		RootDir:        "C:\\Users\\raula\\Desktop\\test_directory\\",
	}
	fileCrawler := crawler.New(app.FileRepo, eventsQueue, app.Logger, crawlerConfig)
	app.Crawler = fileCrawler

	return &app
}
