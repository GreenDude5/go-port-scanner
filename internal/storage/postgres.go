package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewConnection(user, password, host, dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, dbName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed connection config: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}
	log.Println("Connected to PostgreSQL database successfully")
	return db, nil
}

func CreateSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS scan_results (
		id SERIAL PRIMARY KEY,
		port INT NOT NULL,
		status VARCHAR(10) NOT NULL,
		scanned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("Failed to create schema: %w", err)
	}
	return nil
}

func SaveResult(db *sql.DB, port int, status string) error {
	query := `
	INSERT INTO scan_results (port, status)
	VALUES ($1, $2);
	`
	_, err := db.Exec(query, port, status)
	if err != nil {
		return fmt.Errorf("Failed to save result: %w", err)
	}
	return nil
}
