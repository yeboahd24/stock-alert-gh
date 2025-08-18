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
			createIndexesPostgres,
		}
	default: // sqlite
		migrations = []string{
			createUsersTable,
			createUserPreferencesTable,
			createAlertsTable,
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
`

// PostgreSQL-specific table definitions
const createUsersTablePostgres = `
CREATE TABLE IF NOT EXISTS users (
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
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'user_preferences') THEN
        CREATE TABLE user_preferences (
            id TEXT PRIMARY KEY,
            user_id TEXT NOT NULL,
            email_notifications BOOLEAN DEFAULT TRUE,
            push_notifications BOOLEAN DEFAULT TRUE,
            notification_frequency TEXT DEFAULT 'immediate',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        ALTER TABLE user_preferences 
        ADD CONSTRAINT user_preferences_user_id_fkey 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;`

const createAlertsTablePostgres = `
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'alerts') THEN
        CREATE TABLE alerts (
            id TEXT PRIMARY KEY,
            user_id TEXT NOT NULL,
            stock_symbol TEXT NOT NULL,
            stock_name TEXT NOT NULL,
            alert_type TEXT NOT NULL,
            threshold_price REAL,
            current_price REAL,
            status TEXT DEFAULT 'active',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            triggered_at TIMESTAMP
        );
        
        ALTER TABLE alerts 
        ADD CONSTRAINT alerts_user_id_fkey 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;`

const createIndexesPostgres = `
CREATE INDEX IF NOT EXISTS idx_alerts_user_id ON alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status);
CREATE INDEX IF NOT EXISTS idx_alerts_stock_symbol ON alerts(stock_symbol);
CREATE INDEX IF NOT EXISTS idx_alerts_alert_type ON alerts(alert_type);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
`