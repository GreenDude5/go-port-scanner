package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/GreenDude5/go-port-scanner/internal/storage"
)

func StartServer(db *sql.DB) {
	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		results, err := storage.GetResults(db)
		if err != nil {
			http.Error(w, "Failed to get results", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, "Failed to encode results", http.StatusInternalServerError)
		}
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
