package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jamesnetherton/m3u"
)

const (
	// Port is the default port used
	Port = "8080"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = Port
	}
	log.Printf("Starting server on port %v", port)

	http.HandleFunc("/", thing)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func thing(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	path := q.Get("m3u")
	if path == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing 'm3u' parameter"))
		return
	}
	playlist, err := m3u.Parse(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to parse m3u playlist from %s with error: %v", path, err.Error())))
		return
	}
	for i, track := range playlist.Tracks {
		log.Printf("Track name: %v", track.Name)
		log.Printf("Track length: %v", track.Length)
		log.Printf("Track URI: %v", track.URI)
		log.Printf("Track Tags:")
		for i := range track.Tags {
			log.Printf(" - %s=>%s", track.Tags[i].Name, track.Tags[i].Value)
		}
		log.Printf("----------")

		select {
		case <-r.Context().Done():
			log.Printf("Terminated by client")
			return
		default:

			resp, err := http.Get(track.URI)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Unable to download part %d with URI '%s': %s", i, track.URI, err.Error())))
				return
			}

			io.CopyBuffer(w, resp.Body, make([]byte, 100))
		}
	}
}
