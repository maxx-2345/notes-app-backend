package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func respondWithCtx(ctx context.Context, w http.ResponseWriter, r any, httpCode int) {
	if httpCode >= 400 {
		if err, ok := r.(error); ok {
			log.Printf("API Error: %v, status: %d", err, httpCode)
		} else if msg, ok := r.(string); ok {
			log.Printf("API Warning: %s, status: %d", msg, httpCode)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(r)
}

// respond writes a JSON response and logs 4xx/5xx using the request context (for request_id).
func respond(ctx context.Context, w http.ResponseWriter, body any, httpCode int) {
	respondWithCtx(ctx, w, body, httpCode)
}
