package main

import (
	"MyFileExporer/backend/internal/repo/database"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type SearchPayload struct {
	FileName   string
	WordList   []string
	Extensions []string
}

// searchHandler handles the search query sent from the client
func (app *application) searchHandler(w http.ResponseWriter, r *http.Request) {
	// read the payload
	qs := r.URL.Query()

	fileName := qs.Get("file_name")
	wordList := qs["word_list"]
	extensions := qs["extensions"]

	fmt.Printf("file_name=%v\n", fileName)
	fmt.Printf("word_list=%v\n", wordList)
	fmt.Printf("extensions=%v\n", extensions)

	// process the payload
	var fileSearchRequest = database.FileSearchRequest{}

	if fileName != "" {
		fileSearchRequest.Name = &fileName
	}

	if len(wordList) > 0 {
		fileSearchRequest.Words = &wordList
	}

	if len(extensions) > 0 {
		fileSearchRequest.Extension = &extensions
	}

	files, err := app.dbRepo.Files.Search(context.Background(), fileSearchRequest)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if len(files) > 0 {
		var parts []string

		if fileName != "" {
			parts = append(parts, fileName)
		}
		if len(wordList) > 0 {
			parts = append(parts, wordList...)
		}
		if len(extensions) > 0 {
			// Optionally add a prefix to distinguish extensions (e.g., ext:pdf)
			for _, ext := range extensions {
				parts = append(parts, "ext:"+ext)
			}
		}

		query := strings.Join(parts, " ")

		err := app.qdrantRepo.StoreQuery(query)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
	}

	err = jsonResponse(w, http.StatusOK, files)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
