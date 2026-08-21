package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kishanghosh090/GO-MONOLITH/internal/config"
)

func main() {
	fmt.Println(os.Args[1])
	if len(os.Args) < 2 {
		log.Fatal("usage: make migrate <up | down>")
	}
	cfg := config.MustLoad()
	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseUrl,
	)

	if err != nil {
		log.Fatalf("migrate.New: %v", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
