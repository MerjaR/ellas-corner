package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initialises the SQLite database connection
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Database connected successfully")
}

// RunMigrations reads SQL files from the migrations folder and executes them
func RunMigrations() {
	migrationFile := "migrations/create_tables.sql"

	// Log to ensure the migration file is being read
	log.Println("Running migrations from:", migrationFile)

	sqlBytes, err := os.ReadFile(migrationFile)
	if err != nil {
		log.Fatalf("Error reading migration file: %v", err)
	}

	//log.Println("Migration file content:", string(sqlBytes)) // Log the SQL content for debugging

	_, err = DB.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("Error executing migration: %v", err)
	}

	log.Println("Migrations executed successfully")
}
