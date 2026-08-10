package route

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kikudesuyo/room-finder/api/handler"
	v1 "github.com/kikudesuyo/room-finder/api/handler/v1"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", v1.HandleHealth)
		r.Post("/search-profiles", handleFunc(v1.HandleCreateSearchProfile))
		r.Post("/search-profiles/{id}/properties", handleFunc(v1.HandleSaveProperty))
	})

	return r
}

func handleFunc(processFn handler.ProcessFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.HandleRequestAndResponse(r, w, processFn)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigin := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
