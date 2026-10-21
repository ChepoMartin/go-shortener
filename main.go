package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Link struct {
	ShortCode string `json:"short_code"`
	Original  string `json:"original_url"`
	Clicks    int    `json:"clicks"`
}

var (
	store = make(map[string]*Link)
	mu    sync.RWMutex
)

func main() {
	store["hireme"] = &Link{
		ShortCode: "hireme",
		Original:  "https://go.dev",
		Clicks:    0,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /shorten", handleShorten)
	mux.HandleFunc("GET /{code}", handleRedirect)
	mux.HandleFunc("GET /stats/{code}", handleStats)

	fmt.Println("🚀 URL Shortener API is live!")
	fmt.Println("👉 1. Test the redirect: http://localhost:8080/hireme")
	fmt.Println("👉 2. Check click stats: http://localhost:8080/stats/hireme")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	var input Link
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	store[input.ShortCode] = &Link{
		ShortCode: input.ShortCode,
		Original:  input.Original,
		Clicks:    0,
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(store[input.ShortCode])
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	mu.Lock()
	link, exists := store[code]
	if exists {
		link.Clicks++
	}
	mu.Unlock()

	if !exists {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link.Original, http.StatusFound)
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	mu.RLock()
	link, exists := store[code]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}
