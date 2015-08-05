package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

type appHandler func(http.ResponseWriter, *http.Request) (int, error)

// Our appHandler type will now satisify http.Handler
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := fn(w, r)
	log.Printf("Request %v", r)
	if err != nil {
		// We could also log our errors centrally:
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

func myHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("Hello mundo!"))
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
		Handler:      appHandler(myHandler),
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
