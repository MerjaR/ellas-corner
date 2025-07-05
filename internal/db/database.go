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
	// Try common relative path
	pathsToTry := []string{
		"migrations/create_tables.sql",       // when run from root
		"../migrations/create_tables.sql",    // from inside /internal
		"../../migrations/create_tables.sql", // from inside /internal/handlers or repository
	}

	var sqlBytes []byte
	var err error

	for _, path := range pathsToTry {
		sqlBytes, err = os.ReadFile(path)
		if err == nil {
			goto Execute
		}
	}

	return fmt.Errorf("error reading migration file after trying paths: %w", err)

Execute:
	_, err = db.Conn.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("error executing migration: %w", err)
	}

	return nil
}
