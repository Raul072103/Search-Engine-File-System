package main

import "net/http"

// spellCollector suggests queries if there are some misspelled words
func (app *application) spellCollector(w http.ResponseWriter, r *http.Request) {
	// read the payload
	//qs := r.URL.Query()
	//
	//fileName := qs.Get("file_name")
	//wordList := qs["word_list"]
	//extensions := qs["extensions"]
}
