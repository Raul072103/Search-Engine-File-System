package main

import (
	"errors"
	"fmt"
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

	fmt.Println("suggestions", suggestions)

	err = writeJSON(w, http.StatusOK, suggestions)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
