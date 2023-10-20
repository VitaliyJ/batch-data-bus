package http

import "net/http"

// NotFound represents common "not found" response
func NotFound(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`{"error": "not found"}`))
}
