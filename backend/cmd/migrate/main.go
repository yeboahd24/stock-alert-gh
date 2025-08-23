package main

import (
	"database/sql"
	"fmt"
	"log"

	"shares-alert-backend/internal/config"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var db *sql.DB
	switch cfg.Database.Type {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode)
		db, err = sql.Open("postgres", dsn)
	case "sqlite":
		db, err = sql.Open("sqlite3", cfg.Database.FilePath+"?_foreign_keys=on")
	default:
		log.Fatalf("Unsupported database type: %s", cfg.Database.Type)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Adding dividend yield columns...")

	var migrations []string
	if cfg.Database.Type == "postgres" {
		migrations = []string{
			`ALTER TABLE shares_alert_alerts 
			 ADD COLUMN IF NOT EXISTS threshold_yield DECIMAL(5,2),
			 ADD COLUMN IF NOT EXISTS current_yield DECIMAL(5,2),
			 ADD COLUMN IF NOT EXISTS target_yield DECIMAL(5,2),
			 ADD COLUMN IF NOT EXISTS yield_change_threshold DECIMAL(5,2),
			 ADD COLUMN IF NOT EXISTS last_yield DECIMAL(5,2)`,
		}
	} else {
		migrations = []string{
			`ALTER TABLE alerts ADD COLUMN threshold_yield DECIMAL(5,2)`,
			`ALTER TABLE alerts ADD COLUMN current_yield DECIMAL(5,2)`,
			`ALTER TABLE alerts ADD COLUMN target_yield DECIMAL(5,2)`,
			`ALTER TABLE alerts ADD COLUMN yield_change_threshold DECIMAL(5,2)`,
			`ALTER TABLE alerts ADD COLUMN last_yield DECIMAL(5,2)`,
		}
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			log.Printf("Migration %d failed (may already exist): %v", i+1, err)
		} else {
			log.Printf("Migration %d completed successfully", i+1)
		}
	}

	log.Println("Dividend yield migrations completed!")
}