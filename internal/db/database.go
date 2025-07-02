package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Conn *sql.DB
}

// InitDB initialises the SQLite database connection
func InitDB(dataSourceName string) *Database {
	conn, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Database connected successfully")
	return &Database{Conn: conn}
}

// RunMigrations executes migrations using the current DB connection
func (db *Database) RunMigrations() {
	migrationFile := "migrations/create_tables.sql"

	// Log to ensure the migration file is being read
	log.Println("Running migrations from:", migrationFile)

	sqlBytes, err := os.ReadFile(migrationFile)
	if err != nil {
		log.Fatalf("Error reading migration file: %v", err)
	}

	//log.Println("Migration file content:", string(sqlBytes)) // Log the SQL content for debugging

	_, err = db.Conn.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("Error executing migration: %v", err)
	}

	log.Println("Migrations executed successfully")
}
