package main

import (
	"MyFileExporer/common/env"
	"MyFileExporer/common/logger"
	"fmt"
	"go.uber.org/zap"
)

const (
	version = "0.0.0"
)

func main() {
	dbAddress := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.GetString("POSTGRES_USER", "admin"),
		env.GetString("POSTGRES_PASSWORD", "admin_password"),
		env.GetString("POSTGRES_IP", "localhost"),
		env.GetString("POSTGRES_PORT", "5434"),
		env.GetString("POSTGRES_DB", "companies"),
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

	app := &application{
		config: cfg,
		logger: zapLogger,
	}

	// Database
	//database, err := db.New(
	//	cfg.db.addr,
	//	cfg.db.maxOpenConns,
	//	cfg.db.maxIdleConns,
	//	cfg.db.maxIdleTime)
	//defer database.Close()
	//
	//if err != nil {
	//	app.logger.Fatal("database error", zap.Error(err))
	//}

	// Metrics collected
	//expvar.NewString("version").Set(version)
	//expvar.Publish("database", expvar.Func(func() any {
	//	return database.Stats()
	//}))
	//expvar.Publish("goroutines", expvar.Func(func() any {
	//	return runtime.NumGoroutine()
	//}))

	mux := app.mount()
	app.logger.Fatal("server error", zap.Error(app.run(mux)))
}
