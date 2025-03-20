package api

import (
	"MyFileExporer/common/logger"
	"go.uber.org/zap"
)

type application struct {
	config config
	logger *zap.Logger
}

type config struct {
}

func main() {
	var app application
	app.logger = logger.InitLogger("./backend.log")
}
