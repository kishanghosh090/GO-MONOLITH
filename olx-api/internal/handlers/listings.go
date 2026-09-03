package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// listing strtucture
type listing struct {
	id          string
	title       string
	description string
	price       string
	city        string
	created_at  time.Time
}
type ListingHandler struct {
	db *sql.DB
}

func NewListingHandler(db *sql.DB) *ListingHandler {
	return &ListingHandler{
		db: db,
	}
}

func (lh *ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	// req scoped context
	ctx := r.Context()
	rows, err := lh.db.QueryContext(
		ctx,
		`SELECT id,title, description,price , city, created_at, pg_sleep(20) FROM listings
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
func (lh *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	slog.Debug("delete failed", "listing_id", id)
	slog.Info("starting query", "listing_id", id)
	slog.Warn("warn log", "listing_id", id)

	_, err := lh.db.ExecContext(
		ctx,
		`DELETE FROM listings WHERE id = $1`, id)

	if err != nil {
		log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		log.Error("delete failed", "listing_id", id, "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	slog.Info("record deleted", "listing_id", id)
	w.WriteHeader(http.StatusNoContent)
}
