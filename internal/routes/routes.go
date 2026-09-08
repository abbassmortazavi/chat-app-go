package routes

import (
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handelHealth)
	return mux
}
