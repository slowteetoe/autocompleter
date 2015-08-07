package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"slowteetoe.com/autocompleter/wordstore"
	"time"
)

type AutoCompleter interface {
	PrefixSearch(searchTerm string) []string
}

var autoCompleter AutoCompleter

func init() {
	autoCompleter = new(wordstore.Wordstore)
}

type appHandler func(http.ResponseWriter, *http.Request) (int, error)

// Our appHandler type will now satisify http.Handler
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := fn(w, r)
	// log.Printf("Request %v", r)
	if err != nil {
		log.Printf("HTTP %d: %v", err)
		switch status {
		// We can have cases as granular as we like, if we wanted to
		// return custom errors for specific status codes.
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			// Catch any other errors we haven't explicitly handled
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

type autoSuggestionResponse struct {
	Q           string   `json:"q"`
	Suggestions []string `json:"suggestions"`
}

func suggestionsHandler(w http.ResponseWriter, r *http.Request) (int, error) {

	q := r.FormValue("q")

	response := autoSuggestionResponse{Q: q, Suggestions: autoCompleter.PrefixSearch(q)}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
	return 200, nil
}

func main() {
	// since we're deploying to Heroku, we'll need to grab the port from $PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("$PORT was unset, defaulting to %v", port)
	}
	s := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	http.HandleFunc("/favicon.ico", ServeFileHandler)
	http.Handle("/", appHandler(suggestionsHandler))
	log.Fatal(s.ListenAndServe())
}

func ServeFileHandler(res http.ResponseWriter, req *http.Request) {
	fname := path.Base(req.URL.Path)
	http.ServeFile(res, req, "./static/"+fname)
}
