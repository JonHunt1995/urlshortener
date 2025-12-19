package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

var (
	shortcutToURLCache = make(map[string]string)
	cacheMutex         sync.RWMutex
)

type kvPairForDB struct {
	Key string
	Val string
}

func addShortcutFromURL(w http.ResponseWriter, r *http.Request) {
	req := kvPairForDB{}
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cacheMutex.Lock()
	shortcutToURLCache[req.Key] = req.Val
	log.Printf("%v\n", shortcutToURLCache)
	cacheMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func sendToURLFromShortcut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cacheMutex.Lock()
	val, ok := shortcutToURLCache[id]
	if !ok {
		http.NotFound(w, r)
		return
	}
	cacheMutex.Unlock()
	log.Printf("redirecting to %v\n", val)
	http.Redirect(w, r, val, http.StatusFound)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{id}", sendToURLFromShortcut)
	mux.HandleFunc("POST /{$}", addShortcutFromURL)

	log.Print("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
