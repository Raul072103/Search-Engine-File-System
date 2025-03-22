package main

import (
	"MyFileExporer/backend/internal/repo/database"
	"context"
	"net/http"
)

type SearchPayload struct {
	FileName   string   `json:"file_name"`
	WordList   []string `json:"word_list"`
	Extensions []string `json:"extensions"`
}

// searchHandler handles the search query sent from the client
func (app *application) searchHandler(w http.ResponseWriter, r *http.Request) {
	// read the payload
	var searchPayload SearchPayload
	err := readJSON(w, r, &searchPayload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// process the payload
	fileSearchRequest := database.FileSearchRequest{
		Words:     &searchPayload.WordList,
		Extension: &searchPayload.Extensions,
		Name:      &searchPayload.FileName,
	}

	files, err := app.dbRepo.Files.Search(context.Background(), fileSearchRequest)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = jsonResponse(w, http.StatusOK, files)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
