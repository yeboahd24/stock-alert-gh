package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	"shares-alert-backend/internal/cache"
	"shares-alert-backend/internal/config"
	"shares-alert-backend/internal/database"
	"shares-alert-backend/internal/handlers"
	"shares-alert-backend/internal/repository"
	"shares-alert-backend/internal/services"
)

type App struct {
	config       *config.Config
	db           *database.DB
	router       *chi.Mux
	alertService *services.AlertService
}

func New(cfg *config.Config) (*App, error) {
	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(&cache.CacheConfig{
		Host:     cfg.Cache.Host,
		Port:     cfg.Cache.Port,
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
		Enabled:  cfg.Cache.Enabled,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis cache: %w", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	alertRepo := repository.NewAlertRepository(db.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo, &cfg.Auth)
	emailService := services.NewEmailService(&cfg.Email)
	stockCacheTTL := time.Duration(cfg.Cache.StockCacheTTL) * time.Minute
	stockService := services.NewStockService(&cfg.External, redisCache, stockCacheTTL)
	alertService := services.NewAlertService(alertRepo, userRepo, stockService, emailService)
	cacheService := services.NewCacheService(redisCache)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	stockHandler := handlers.NewStockHandler(stockService)
	alertHandler := handlers.NewAlertHandler(alertService)
	userHandler := handlers.NewUserHandler(userRepo)
	cacheHandler := handlers.NewCacheHandler(cacheService, stockService)

	// Setup router
	router := setupRouter(cfg, authHandler, stockHandler, alertHandler, userHandler, cacheHandler)

	app := &App{
		config:       cfg,
		db:           db,
		router:       router,
		alertService: alertService,
	}

	// Start alert monitoring in background
	go app.alertService.StartMonitoring()

	return app, nil
}

func (a *App) Start(addr string) error {
	log.Printf("Server starting on %s", addr)
	log.Printf("Health check: http://localhost%s/api/v1/health", addr)
	return http.ListenAndServe(addr, a.router)
}

func setupRouter(
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	stockHandler *handlers.StockHandler,
	alertHandler *handlers.AlertHandler,
	userHandler *handlers.UserHandler,
	cacheHandler *handlers.CacheHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(time.Duration(cfg.Server.RequestTimeout) * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Server.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{
				"status":    "ok",
				"timestamp": time.Now().Format(time.RFC3339),
				"version":   "2.0.0",
			})
		})

		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Get("/google", authHandler.GetGoogleAuthURL)
			r.Post("/google/callback", authHandler.GoogleCallback)
			r.With(authHandler.AuthMiddleware).Post("/logout", authHandler.Logout)
			r.With(authHandler.AuthMiddleware).Get("/profile", authHandler.GetProfile)
		})

		// Stock routes (public, but can be enhanced with auth)
		r.Route("/stocks", func(r chi.Router) {
			r.Use(authHandler.OptionalAuthMiddleware)
			r.Get("/", stockHandler.GetAllStocks)
			r.Get("/{symbol}", stockHandler.GetStock)
			r.Get("/{symbol}/details", stockHandler.GetStockDetails)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authHandler.AuthMiddleware)

			// Alert routes
			r.Route("/alerts", func(r chi.Router) {
				r.Get("/", alertHandler.GetAlerts)
				r.Post("/", alertHandler.CreateAlert)
				r.Get("/{id}", alertHandler.GetAlert)
				r.Put("/{id}", alertHandler.UpdateAlert)
				r.Delete("/{id}", alertHandler.DeleteAlert)
			})

			// User routes
			r.Route("/user", func(r chi.Router) {
				r.Get("/preferences", userHandler.GetPreferences)
				r.Put("/preferences", userHandler.UpdatePreferences)
			})

			// Cache management routes (admin only in production)
			r.Route("/cache", func(r chi.Router) {
				r.Get("/stats", cacheHandler.GetCacheStats)
				r.Post("/invalidate", cacheHandler.InvalidateCache)
				r.Post("/warmup", cacheHandler.WarmupCache)
			})
		})
	})

	return r
}