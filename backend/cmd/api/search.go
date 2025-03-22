package main

import (
	"net/http"
)

type SearchPayload struct {
	FileName         string   `json:"file_name"`
	WordList         []string `json:"word_list"`
	ExtensionInclude string   `json:"extension_include"`
	ExtensionExclude string   `json:"extension_exclude"`
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

}
