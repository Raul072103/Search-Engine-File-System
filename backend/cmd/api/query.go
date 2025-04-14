package main

import (
	"errors"
	"net/http"
)

var QuerySuggestionsLimit = uint64(5)

func (app *application) querySuggestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		app.badRequestResponse(w, r, errors.New("no query parameter"))
		return
	}
	suggestions, err := app.qdrantRepo.SuggestSimilarQueries(query, &QuerySuggestionsLimit)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, suggestions)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
