package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Email    EmailConfig
	External ExternalConfig
	Cache    CacheConfig
}

type ServerConfig struct {
	Port            string
	AllowedOrigins  []string
	RequestTimeout  int
}

type DatabaseConfig struct {
	Type     string // sqlite, postgres, mysql
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	FilePath string // for SQLite
}

type AuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	RedirectURL        string
	JWTSecret          string
	JWTExpirationHours int
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type ExternalConfig struct {
	GSEBaseURL string
	ProxyURL   string
}

type CacheConfig struct {
	URL       string
	Host      string
	Port      string
	Password  string
	DB        int
	Enabled   bool
	StockCacheTTL int // in minutes
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "10000"),
			AllowedOrigins: []string{
				getEnv("FRONTEND_URL", "http://localhost:3000"),
				"http://localhost:3000",
				"http://localhost:5173",
				"https://stock-alert-gh.onrender.com",
			},
			RequestTimeout: getEnvAsInt("REQUEST_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "sqlite"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "shares_alert"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			FilePath: getEnv("DB_FILE_PATH", "./data/shares_alert.db"),
		},
		Auth: AuthConfig{
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:        getEnv("OAUTH_REDIRECT_URL", "http://localhost:5173/"),
			JWTSecret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnv("SMTP_PORT", "587"),
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("FROM_EMAIL", ""),
			FromName:     getEnv("FROM_NAME", "Shares Alert Ghana"),
		},
		External: ExternalConfig{
			GSEBaseURL: getEnv("GSE_BASE_URL", "https://dev.kwayisi.org/apis/gse"),
			ProxyURL:   getEnv("PROXY_URL", "https://api.allorigins.win/raw?url="),
		},
		Cache: CacheConfig{
			URL:           getEnv("REDIS_URL", ""),
			Host:          getEnv("REDIS_HOST", "localhost"),
			Port:          getEnv("REDIS_PORT", "6379"),
			Password:      getEnv("REDIS_PASSWORD", ""),
			DB:            getEnvAsInt("REDIS_DB", 0),
			Enabled:       getEnvAsBool("REDIS_ENABLED", true),
			StockCacheTTL: getEnvAsInt("STOCK_CACHE_TTL_MINUTES", 5),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}