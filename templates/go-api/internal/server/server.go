package server

import (
	"net/http"

	"github.com/{{ .Author }}/{{ .ProjectName }}/internal/handlers"
)

func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)
	return mux
}
