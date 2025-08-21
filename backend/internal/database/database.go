package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"shares-alert-backend/internal/config"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
	config *config.DatabaseConfig
}

func New(cfg *config.DatabaseConfig) (*DB, error) {
	var db *sql.DB
	var err error

	switch cfg.Type {
	case "sqlite":
		// Ensure data directory exists
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
		
		db, err = sql.Open("sqlite3", cfg.FilePath+"?_foreign_keys=on")
		if err != nil {
			return nil, fmt.Errorf("failed to open SQLite database: %w", err)
		}
		
		// Set SQLite pragmas for better performance
		if _, err := db.Exec(`
			PRAGMA journal_mode = WAL;
			PRAGMA synchronous = NORMAL;
			PRAGMA cache_size = 1000;
			PRAGMA foreign_keys = on;
			PRAGMA temp_store = memory;
		`); err != nil {
			return nil, fmt.Errorf("failed to set SQLite pragmas: %w", err)
		}

	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open PostgreSQL database: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbInstance := &DB{
		DB:     db,
		config: cfg,
	}

	// Run migrations
	if err := dbInstance.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Printf("Connected to %s database successfully", cfg.Type)
	return dbInstance, nil
}

func (db *DB) Migrate() error {
	var migrations []string
	
	// Use different migrations based on database type
	switch db.config.Type {
	case "postgres":
		migrations = []string{
			createUsersTablePostgres,
			createUserPreferencesTablePostgres,
			createAlertsTablePostgres,
			createIPOTablePostgres,
			createIndexesPostgres,
		}
	default: // sqlite
		migrations = []string{
			createUsersTable,
			createUserPreferencesTable,
			createAlertsTable,
			createIPOTable,
			createIndexes,
		}
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	name TEXT NOT NULL,
	picture TEXT,
	google_id TEXT UNIQUE NOT NULL,
	email_verified BOOLEAN DEFAULT FALSE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createUserPreferencesTable = `
CREATE TABLE IF NOT EXISTS user_preferences (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	email_notifications BOOLEAN DEFAULT TRUE,
	push_notifications BOOLEAN DEFAULT TRUE,
	notification_frequency TEXT DEFAULT 'immediate',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createAlertsTable = `
CREATE TABLE IF NOT EXISTS alerts (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	stock_symbol TEXT NOT NULL,
	stock_name TEXT NOT NULL,
	alert_type TEXT NOT NULL,
	threshold_price REAL,
	current_price REAL,
	status TEXT DEFAULT 'active',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	triggered_at DATETIME,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_alerts_user_id ON alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status);
CREATE INDEX IF NOT EXISTS idx_alerts_stock_symbol ON alerts(stock_symbol);
CREATE INDEX IF NOT EXISTS idx_alerts_alert_type ON alerts(alert_type);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
CREATE INDEX IF NOT EXISTS idx_ipo_status ON ipo_announcements(status);
CREATE INDEX IF NOT EXISTS idx_ipo_listing_date ON ipo_announcements(listing_date);
CREATE INDEX IF NOT EXISTS idx_ipo_symbol ON ipo_announcements(symbol);
`

// PostgreSQL-specific table definitions
const createUsersTablePostgres = `
CREATE TABLE IF NOT EXISTS shares_alert_users (
	id TEXT PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	name TEXT NOT NULL,
	picture TEXT,
	google_id TEXT UNIQUE NOT NULL,
	email_verified BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createUserPreferencesTablePostgres = `
CREATE TABLE IF NOT EXISTS shares_alert_user_preferences (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES shares_alert_users(id) ON DELETE CASCADE,
	email_notifications BOOLEAN DEFAULT TRUE,
	push_notifications BOOLEAN DEFAULT TRUE,
	notification_frequency TEXT DEFAULT 'immediate',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createAlertsTablePostgres = `
CREATE TABLE IF NOT EXISTS shares_alert_alerts (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES shares_alert_users(id) ON DELETE CASCADE,
	stock_symbol TEXT NOT NULL,
	stock_name TEXT NOT NULL,
	alert_type TEXT NOT NULL,
	threshold_price REAL,
	current_price REAL,
	status TEXT DEFAULT 'active',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	triggered_at TIMESTAMP
);`

const createIPOTable = `
CREATE TABLE IF NOT EXISTS ipo_announcements (
	id TEXT PRIMARY KEY,
	company_name TEXT NOT NULL,
	symbol TEXT NOT NULL UNIQUE,
	sector TEXT,
	offer_price REAL NOT NULL,
	listing_date DATETIME NOT NULL,
	status TEXT NOT NULL DEFAULT 'announced',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createIPOTablePostgres = `
CREATE TABLE IF NOT EXISTS shares_alert_ipo_announcements (
	id TEXT PRIMARY KEY,
	company_name TEXT NOT NULL,
	symbol TEXT NOT NULL UNIQUE,
	sector TEXT,
	offer_price REAL NOT NULL,
	listing_date TIMESTAMP NOT NULL,
	status TEXT NOT NULL DEFAULT 'announced',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createIndexesPostgres = `
CREATE INDEX IF NOT EXISTS idx_shares_alert_alerts_user_id ON shares_alert_alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_shares_alert_alerts_status ON shares_alert_alerts(status);
CREATE INDEX IF NOT EXISTS idx_shares_alert_alerts_stock_symbol ON shares_alert_alerts(stock_symbol);
CREATE INDEX IF NOT EXISTS idx_shares_alert_alerts_alert_type ON shares_alert_alerts(alert_type);
CREATE INDEX IF NOT EXISTS idx_shares_alert_users_email ON shares_alert_users(email);
CREATE INDEX IF NOT EXISTS idx_shares_alert_users_google_id ON shares_alert_users(google_id);
CREATE INDEX IF NOT EXISTS idx_shares_alert_ipo_status ON shares_alert_ipo_announcements(status);
CREATE INDEX IF NOT EXISTS idx_shares_alert_ipo_listing_date ON shares_alert_ipo_announcements(listing_date);
CREATE INDEX IF NOT EXISTS idx_shares_alert_ipo_symbol ON shares_alert_ipo_announcements(symbol);
`