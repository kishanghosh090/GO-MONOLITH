package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type listing struct {
	id          string
	title       string
	description string
	price       string
	city        string
	created_at  time.Time
}

func Listings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(
			`SELECT id,title, description,price , city, created_at FROM listings
			ORDER BY created_at DESC
			LIMIT 100
			`)
		if err != nil {
			log.Printf("query: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		listings := []listing{}
		for rows.Next() {
			var l listing
			if err := rows.Scan(&l.id, &l.title, &l.description, &l.price, &l.city, &l.created_at); err != nil {
				log.Printf("Rows.Err: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			listings = append(listings, l)
		}

		if err := rows.Err(); err != nil {
			log.Printf("rows.err: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(listings)

	}
}
