package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Conn *sql.DB
}

// InitDB initialises the SQLite database connection
func InitDB(dataSourceName string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return &Database{Conn: conn}, nil
}

// RunMigrations executes migrations using the current DB connection
func (db *Database) RunMigrations() error {
	migrationFile := "migrations/create_tables.sql"

	sqlBytes, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("error reading migration file: %w", err)
	}

	_, err = db.Conn.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("error executing migration: %w", err)
	}

	return nil
}
