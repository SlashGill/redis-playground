package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"url-shortener-basic/internal/service"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	shortCode, err := service.Shorten(req.URL)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"short_url": "http://localhost:9999/" + shortCode})
}

func ResolveURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:] // 移除前導 "/"

	originalURL, err := service.Resolve(shortCode)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	log.Println("originalURL:", originalURL)
	http.Redirect(w, r, originalURL, http.StatusFound)
}
