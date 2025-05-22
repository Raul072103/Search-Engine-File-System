package main

import (
	"fmt"
	"net/http"
)

type Suggestion struct {
	FileNameSuggestion  string   `json:"file_name_suggestion"`
	WordListSuggestions []string `json:"word_list_suggestions"`
}

// spellCollector suggests queries if there are some misspelled words
func (app *application) spellCollector(w http.ResponseWriter, r *http.Request) {
	// read the payload
	qs := r.URL.Query()

	fileName := qs.Get("file_name")
	wordList := qs["word_list"]

	suggestion := Suggestion{
		FileNameSuggestion:  "",
		WordListSuggestions: make([]string, 0),
	}

	fmt.Println("AICI")

	if app.spellingCorrector.Initialized() {
		// file name
		if fileName != "" {
			fileNameSuggestion := app.spellingCorrector.Correction(fileName)
			if fileNameSuggestion != "" {
				suggestion.FileNameSuggestion = fileNameSuggestion
			} else {
				suggestion.FileNameSuggestion = fileName
			}
		}

		// words
		for _, word := range wordList {
			wordSuggestion := app.spellingCorrector.Correction(word)
			if wordSuggestion != "" {
				suggestion.WordListSuggestions = append(suggestion.WordListSuggestions, wordSuggestion)
			} else {
				suggestion.WordListSuggestions = append(suggestion.WordListSuggestions, word)
			}
		}

	} else {
		err := writeJSON(w, http.StatusOK, "")
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
	}

	err := writeJSON(w, http.StatusOK, suggestion)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
