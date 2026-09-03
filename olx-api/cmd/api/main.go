package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kishanghosh090/GO-MONOLITH/internal/config"
	"github.com/kishanghosh090/GO-MONOLITH/internal/db"
	"github.com/kishanghosh090/GO-MONOLITH/internal/handlers"
)

func main() {
	// load env
	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.Connect: %v", err)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	lh := handlers.NewListingHandler(db)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

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
