package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase initializes the database (SQLite or PostgreSQL)
func InitDatabase(dbPath string) error {
	var err error
	
	// Check if DATABASE_URL is set (Railway PostgreSQL)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Println("Using PostgreSQL database")
		DB, err = sql.Open("postgres", databaseURL)
	} else {
		log.Println("Using SQLite database")
		DB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
	}
	
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Execute schema
	if err = executeSchema(); err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// executeSchema runs the SQL schema file
func executeSchema() error {
	schemaPath := filepath.Join("database", "schema.sql")
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	schema := string(schemaBytes)
	if _, err := DB.Exec(schema); err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}