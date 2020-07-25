package controller

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rl404/point-system/internal/config"
	"github.com/rl404/point-system/internal/view"
)

// RegisterBaseRoutes registers base routes.
func RegisterBaseRoutes(router *chi.Mux) {
	// Root route for testing.
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		view.RespondWithJSON(w, 200, "root", nil)
	})

	// Ping route also for testing.
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		view.RespondWithJSON(w, 200, "pong", nil)
	})

	// Handle page not found 404.
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		view.RespondWithJSON(w, 404, "page not found", nil)
	}))
}

// GetRoutesV1 to get all v1 routes.
func GetRoutesV1(cfg config.Config) (r http.Handler, err error) {
	router := chi.NewRouter()

	ph, err := newPointHandler(cfg)
	if err != nil {
		return router, err
	}

	router.Get("/point", ph.getPoint)
	router.Post("/point/add", ph.addPoint)
	router.Post("/point/subtract", ph.subtractPoint)

	return router, nil
}
