package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
    DB *sql.DB
}

// Open opens the SQLite database, enables foreign keys, and runs migrations.
func Open(path string, migrationsPath string) (*Database, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Ensure the DB is reachable
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    // Enable foreign key constraints
    _, err = db.Exec("PRAGMA foreign_keys = ON;")
    if err != nil {
        return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
    }

    // Run migrations
    if err := runMigrations(db, migrationsPath); err != nil {
        return nil, fmt.Errorf("failed to run migrations: %w", err)
    }

    return &Database{DB: db}, nil
}

// runMigrations reads and executes the SQL migration file.
func runMigrations(db *sql.DB, path string) error {
    sqlBytes, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read migrations file: %w", err)
    }

    _, err = db.Exec(string(sqlBytes))
    if err != nil {
        return fmt.Errorf("failed to execute migrations: %w", err)
    }

    return nil
}

// Utility function for generating timestamps
func Now() string {
    return time.Now().Format(time.RFC3339)
}

/*package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Connect() error {
    var err error
    DB, err = sql.Open("sqlite3", "./forum.db")
    if err != nil {
        return err
    }

    return DB.Ping()
}

func RunMigrations() error {
    // Read the migrations.sql file
    content, err := os.ReadFile("backend/database/migrations.sql")
    if err != nil {
        return fmt.Errorf("failed to read migrations.sql: %v", err)
    }

    // Convert to string
    queries := string(content)

    // Split into individual statements
    statements := strings.Split(queries, ";")

    // Execute each statement
    for _, stmt := range statements {
        stmt = strings.TrimSpace(stmt)
        if stmt == "" {
            continue
        }

        _, err := DB.Exec(stmt)
        if err != nil {
            log.Printf("Migration failed on statement: %s\nError: %v\n", stmt, err)
            return err
        }
    }

    log.Println("Migrations executed successfully.")
    return nil
}*/