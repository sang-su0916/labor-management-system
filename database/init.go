package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase initializes the database (SQLite or PostgreSQL)
func InitDatabase(dbPath string) error {
	var err error
	
	// Check if DATABASE_URL is set (Render PostgreSQL)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Println("🐘 Using PostgreSQL database")
		log.Printf("📊 Database URL: %s", maskDatabaseURL(databaseURL))
		DB, err = sql.Open("postgres", databaseURL)
		if err != nil {
			return fmt.Errorf("failed to open PostgreSQL connection: %v", err)
		}
	} else {
		log.Println("🗄️ Using SQLite database")
		log.Printf("📊 Database path: %s", dbPath)
		DB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
		if err != nil {
			return fmt.Errorf("failed to open SQLite connection: %v", err)
		}
	}

	// Test the connection with retry
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		log.Printf("⚠️ Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	
	if err != nil {
		return fmt.Errorf("failed to ping database after %d attempts: %v", maxRetries, err)
	}

	log.Println("✅ Database connection established")

	// Execute schema
	if err = executeSchema(); err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	log.Println("✅ Database initialized successfully")
	return nil
}

// maskDatabaseURL masks sensitive information in database URL for logging
func maskDatabaseURL(url string) string {
	// Simple masking - hide password
	if len(url) > 20 {
		return url[:10] + "****" + url[len(url)-10:]
	}
	return "****"
}

// executeSchema runs the SQL schema file
func executeSchema() error {
	var schemaBytes []byte
	var err error
	
	// Check if DATABASE_URL is set (PostgreSQL)
	if os.Getenv("DATABASE_URL") != "" {
		log.Println("Using PostgreSQL - loading postgres_schema.sql")
		// Try PostgreSQL schema paths
		possiblePaths := []string{
			"database/postgres_schema.sql",
			"./database/postgres_schema.sql",
			"postgres_schema.sql",
		}
		
		for _, path := range possiblePaths {
			schemaBytes, err = ioutil.ReadFile(path)
			if err == nil {
				log.Printf("PostgreSQL schema file found at: %s", path)
				break
			}
		}
		
		if err != nil {
			log.Printf("PostgreSQL schema file not found, creating basic tables manually")
			return createPostgreSQLTables()
		}
	} else {
		log.Println("Using SQLite - loading schema.sql")
		// Try SQLite schema paths
		possiblePaths := []string{
			"database/schema.sql",
			"./database/schema.sql",
			"schema.sql",
		}
		
		for _, path := range possiblePaths {
			schemaBytes, err = ioutil.ReadFile(path)
			if err == nil {
				log.Printf("SQLite schema file found at: %s", path)
				break
			}
		}
		
		if err != nil {
			return fmt.Errorf("failed to read SQLite schema file from any location: %v", err)
		}
	}

	schema := string(schemaBytes)
	if _, err := DB.Exec(schema); err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	log.Println("Database schema executed successfully")
	return nil
}

// createPostgreSQLTables creates PostgreSQL tables manually if schema file is not found
func createPostgreSQLTables() error {
	log.Println("Creating PostgreSQL tables manually...")
	
	// Create users table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			role VARCHAR(20) DEFAULT 'employee',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}
	
	// Create employees table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			employee_number VARCHAR(20) UNIQUE NOT NULL,
			department VARCHAR(50) NOT NULL,
			position VARCHAR(50) NOT NULL,
			hire_date DATE NOT NULL,
			salary DECIMAL(12,2) DEFAULT 0,
			phone VARCHAR(20),
			email VARCHAR(100),
			address TEXT,
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create employees table: %v", err)
	}
	
	// Create system_settings table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS system_settings (
			id SERIAL PRIMARY KEY,
			setting_key VARCHAR(100) UNIQUE NOT NULL,
			setting_value TEXT,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create system_settings table: %v", err)
	}
	
	return createInitialData()
}

// createInitialData creates initial admin user and system settings
func createInitialData() error {
	// Create basic admin user if not exists
	_, err := DB.Exec(`
		INSERT INTO users (username, password_hash, email, role)
		VALUES ('admin', '$2a$10$BGuuHyAsIfgXDObMqhNUwOnfY4oK56B50BVx1NoZWL0y9kRmsdYji', 'admin@company.com', 'admin')
		ON CONFLICT (username) DO NOTHING
	`)
	
	if err != nil {
		log.Printf("Note: Could not create initial admin user: %v", err)
	} else {
		log.Println("Initial admin user created successfully")
	}
	
	// Insert basic system settings
	_, err = DB.Exec(`
		INSERT INTO system_settings (setting_key, setting_value, description) VALUES
		('company_name', '테스트 회사', '회사명'),
		('company_address', '서울특별시 강남구 테헤란로 123', '회사 주소'),
		('company_phone', '02-1234-5678', '회사 전화번호')
		ON CONFLICT (setting_key) DO NOTHING
	`)
	
	if err != nil {
		log.Printf("Note: Could not create initial system settings: %v", err)
	} else {
		log.Println("Initial system settings created successfully")
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