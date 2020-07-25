package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rl404/point-system/internal/config"
	"github.com/rl404/point-system/internal/controller"
	"github.com/rs/cors"
)

// startHTTP to start serving HTTP.
func startHTTP(cfg config.Config) error {
	r := chi.NewRouter()

	// Set default recommended go-chi router middlewares.
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.Default().Handler)

	// Register base routes.
	controller.RegisterBaseRoutes(r)

	// Prepare route V1.
	v1Route, err := controller.GetRoutesV1(cfg)
	if err != nil {
		return err
	}

	// Register v1 routes.
	r.Mount("/v1", v1Route)

	fmt.Println("server listen at " + cfg.Port)
	return http.ListenAndServe(cfg.Port, r)
}

func main() {
	err := startHTTP(config.GetConfig())
	if err != nil {
		log.Fatal("error starting HTTP", " - ", err.Error())
		return
	}
}
