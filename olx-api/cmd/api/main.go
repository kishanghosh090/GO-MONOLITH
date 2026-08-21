package main

import (
	"log"
	"net/http"
	"time"

	"github.com/kishanghosh090/GO-MONOLITH/internal/config"
	"github.com/kishanghosh090/GO-MONOLITH/internal/db"
	"github.com/kishanghosh090/GO-MONOLITH/internal/handlers"
)

func main() {
	// load env
	cfg := config.MustLoad()
	_, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.Connect: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", handlers.Listings)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("Server is listing on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
