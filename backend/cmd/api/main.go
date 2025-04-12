package main

import (
	"MyFileExporer/backend/internal/db"
	"MyFileExporer/backend/internal/repo/database"
	"MyFileExporer/common/env"
	"MyFileExporer/common/logger"
	"expvar"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"runtime"
)

const (
	version = "0.0.0"
)

func main() {
	// Env file setup
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("Error loading .env file", zap.Error(err))
	}

	dbAddress := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.GetString("POSTGRES_USER", "admin"),
		env.GetString("POSTGRES_PASSWORD", "admin_password"),
		env.GetString("POSTGRES_IP", "localhost"),
		env.GetString("POSTGRES_PORT", "5434"),
		env.GetString("POSTGRES_DB", "file_system_database"),
	)

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         dbAddress,
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:         env.GetString("DEV", "development"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5174"),
	}

	zapLogger := logger.InitLogger("backend.log")

	// Database
	pgDb, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime)
	defer pgDb.Close()

	if err != nil {
		zapLogger.Fatal("db error", zap.Error(err))
	}

	dbRepo := database.NewRepo(pgDb)

	app := &application{
		config: cfg,
		logger: zapLogger,
		dbRepo: dbRepo,
	}

	// Metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("db", expvar.Func(func() any {
		return pgDb.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	app.logger.Fatal("server error", zap.Error(app.run(mux)))

	err = app.run(mux)
	if err != nil {
		app.logger.Fatal("server error", zap.Error(err))
	}
}
