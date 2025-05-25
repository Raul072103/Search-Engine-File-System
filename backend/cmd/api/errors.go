package main

import (
	"go.uber.org/zap"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("internal error", zap.String("method", r.Method),
		zap.String("path", r.URL.Path), zap.String("error", err.Error()))

	_ = writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("bad request", zap.String("method", r.Method),
		zap.String("path", r.URL.Path), zap.String("error", err.Error()))

	_ = writeJSONError(w, http.StatusBadRequest, err.Error())
}
