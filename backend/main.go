package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// Stock data structures based on the Ghana Stock Exchange API
type StockLive struct {
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Change float64 `json:"change"`
	Volume int64   `json:"volume"`
}

type Company struct {
	Address   string   `json:"address"`
	Directors []string `json:"directors"`
	Email     string   `json:"email"`
	Facsimile *string  `json:"facsimile"`
	Industry  string   `json:"industry"`
	Name      string   `json:"name"`
	Sector    string   `json:"sector"`
	Telephone string   `json:"telephone"`
	Website   string   `json:"website"`
}

type StockEquity struct {
	Capital float64  `json:"capital"`
	Company Company  `json:"company"`
	DPS     *float64 `json:"dps"`
	EPS     *float64 `json:"eps"`
	Name    string   `json:"name"`
	Price   float64  `json:"price"`
	Shares  int64    `json:"shares"`
}

// Enhanced stock data for our application
type EnhancedStock struct {
	Symbol           string    `json:"symbol"`
	Name             string    `json:"name"`
	CurrentPrice     float64   `json:"currentPrice"`
	PreviousClose    float64   `json:"previousClose"`
	Change           float64   `json:"change"`
	ChangePercent    float64   `json:"changePercent"`
	Volume           int64     `json:"volume"`
	LastUpdated      time.Time `json:"lastUpdated"`
	MarketCap        *float64  `json:"marketCap,omitempty"`
	Sector           *string   `json:"sector,omitempty"`
	Industry         *string   `json:"industry,omitempty"`
}

// Detailed stock data including company information
type DetailedStock struct {
	Symbol           string    `json:"symbol"`
	Name             string    `json:"name"`
	CurrentPrice     float64   `json:"currentPrice"`
	PreviousClose    float64   `json:"previousClose"`
	Change           float64   `json:"change"`
	ChangePercent    float64   `json:"changePercent"`
	Volume           int64     `json:"volume"`
	LastUpdated      time.Time `json:"lastUpdated"`
	MarketCap        float64   `json:"marketCap"`
	Shares           int64     `json:"shares"`
	Sector           string    `json:"sector"`
	Industry         string    `json:"industry"`
	DPS              *float64  `json:"dps"`
	EPS              *float64  `json:"eps"`
	Company          Company   `json:"company"`
}

// Alert structures
type Alert struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	StockSymbol   string    `json:"stockSymbol"`
	StockName     string    `json:"stockName"`
	AlertType     string    `json:"alertType"`
	ThresholdPrice *float64 `json:"thresholdPrice,omitempty"`
	CurrentPrice  *float64  `json:"currentPrice,omitempty"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Alert request/response structures
type CreateAlertRequest struct {
	StockSymbol    string   `json:"stockSymbol"`
	StockName      string   `json:"stockName"`
	AlertType      string   `json:"alertType"`
	ThresholdPrice *float64 `json:"thresholdPrice,omitempty"`
}

type UpdateAlertRequest struct {
	AlertType      *string  `json:"alertType,omitempty"`
	ThresholdPrice *float64 `json:"thresholdPrice,omitempty"`
	Status         *string  `json:"status,omitempty"`
}

// In-memory storage (in production, use a database)
var alerts []Alert
var alertCounter int

const GSE_BASE_URL = "https://dev.kwayisi.org/apis/gse"

func main() {
	// Initialize sample data
	initSampleAlerts()
	
	// Start alert monitoring in background
	go startAlertMonitoring()
	
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "https://stock-alert-gh.onrender.com", "https://stock-alert-gh-backend.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Stock routes
		r.Route("/stocks", func(r chi.Router) {
			r.Get("/", getAllStocks)
			r.Get("/{symbol}", getStock)
			r.Get("/{symbol}/details", getStockDetails)
		})

		// Alert routes
		r.Route("/alerts", func(r chi.Router) {
			r.Get("/", getAlerts)
			r.Post("/", createAlert)
			r.Get("/{id}", getAlert)
			r.Put("/{id}", updateAlert)
			r.Delete("/{id}", deleteAlert)
		})

		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{"status": "ok", "timestamp": time.Now().Format(time.RFC3339)})
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Ghana Stock Exchange API Backend\n")
	fmt.Printf("Health check: http://localhost:%s/api/v1/health\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Stock handlers
func getAllStocks(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(GSE_BASE_URL + "/live")
	if err != nil {
		http.Error(w, "Failed to fetch stock data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var stocks []StockLive
	if err := json.NewDecoder(resp.Body).Decode(&stocks); err != nil {
		http.Error(w, "Failed to parse stock data", http.StatusInternalServerError)
		return
	}

	// Convert to enhanced format
	enhancedStocks := make([]EnhancedStock, len(stocks))
	for i, stock := range stocks {
		changePercent := 0.0
		if stock.Price > 0 {
			previousClose := stock.Price - stock.Change
			if previousClose > 0 {
				changePercent = (stock.Change / previousClose) * 100
			}
		}

		enhancedStocks[i] = EnhancedStock{
			Symbol:        stock.Name,
			Name:          stock.Name,
			CurrentPrice:  stock.Price,
			PreviousClose: stock.Price - stock.Change,
			Change:        stock.Change,
			ChangePercent: changePercent,
			Volume:        stock.Volume,
			LastUpdated:   time.Now(),
		}
	}

	render.JSON(w, r, enhancedStocks)
}

func getStock(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	symbol = strings.ToUpper(symbol)

	resp, err := http.Get(GSE_BASE_URL + "/live/" + symbol)
	if err != nil {
		http.Error(w, "Failed to fetch stock data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	var stock StockLive
	if err := json.NewDecoder(resp.Body).Decode(&stock); err != nil {
		http.Error(w, "Failed to parse stock data", http.StatusInternalServerError)
		return
	}

	// Convert to enhanced format
	changePercent := 0.0
	if stock.Price > 0 {
		previousClose := stock.Price - stock.Change
		if previousClose > 0 {
			changePercent = (stock.Change / previousClose) * 100
		}
	}

	enhancedStock := EnhancedStock{
		Symbol:        stock.Name,
		Name:          stock.Name,
		CurrentPrice:  stock.Price,
		PreviousClose: stock.Price - stock.Change,
		Change:        stock.Change,
		ChangePercent: changePercent,
		Volume:        stock.Volume,
		LastUpdated:   time.Now(),
	}

	render.JSON(w, r, enhancedStock)
}

func getStockDetails(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	symbol = strings.ToUpper(symbol)

	// Get equity details
	resp, err := http.Get(GSE_BASE_URL + "/equities/" + symbol)
	if err != nil {
		http.Error(w, "Failed to fetch stock details", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	var equity StockEquity
	if err := json.NewDecoder(resp.Body).Decode(&equity); err != nil {
		http.Error(w, "Failed to parse stock details", http.StatusInternalServerError)
		return
	}

	// Get live data for current price, change, and volume
	var liveStock StockLive
	var currentPrice float64 = equity.Price
	var change float64 = 0
	var volume int64 = 0
	var changePercent float64 = 0

	liveResp, err := http.Get(GSE_BASE_URL + "/live/" + symbol)
	if err == nil {
		defer liveResp.Body.Close()
		if liveResp.StatusCode == 200 {
			if json.NewDecoder(liveResp.Body).Decode(&liveStock) == nil {
				currentPrice = liveStock.Price
				change = liveStock.Change
				volume = liveStock.Volume
				
				if currentPrice > 0 {
					previousClose := currentPrice - change
					if previousClose > 0 {
						changePercent = (change / previousClose) * 100
					}
				}
			}
		}
	}

	detailedStock := DetailedStock{
		Symbol:        equity.Name,
		Name:          equity.Company.Name,
		CurrentPrice:  currentPrice,
		PreviousClose: currentPrice - change,
		Change:        change,
		ChangePercent: changePercent,
		Volume:        volume,
		LastUpdated:   time.Now(),
		MarketCap:     equity.Capital,
		Shares:        equity.Shares,
		Sector:        equity.Company.Sector,
		Industry:      equity.Company.Industry,
		DPS:           equity.DPS,
		EPS:           equity.EPS,
		Company:       equity.Company,
	}

	render.JSON(w, r, detailedStock)
}

// Alert handlers
func getAlerts(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	status := r.URL.Query().Get("status")
	
	var filteredAlerts []Alert
	for _, alert := range alerts {
		// Filter by user ID if provided
		if userID != "" && alert.UserID != userID {
			continue
		}
		// Filter by status if provided
		if status != "" && alert.Status != status {
			continue
		}
		filteredAlerts = append(filteredAlerts, alert)
	}
	
	render.JSON(w, r, filteredAlerts)
}

func createAlert(w http.ResponseWriter, r *http.Request) {
	var req CreateAlertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.StockSymbol == "" || req.AlertType == "" {
		http.Error(w, "stockSymbol and alertType are required", http.StatusBadRequest)
		return
	}

	// Validate alert type
	validAlertTypes := map[string]bool{
		"price_threshold":        true,
		"ipo_alert":             true,
		"dividend_announcement": true,
	}
	if !validAlertTypes[req.AlertType] {
		http.Error(w, "Invalid alert type", http.StatusBadRequest)
		return
	}

	// For price threshold alerts, threshold price is required
	if req.AlertType == "price_threshold" && req.ThresholdPrice == nil {
		http.Error(w, "thresholdPrice is required for price_threshold alerts", http.StatusBadRequest)
		return
	}

	// Get current stock price
	var currentPrice *float64
	if stockResp, err := http.Get(GSE_BASE_URL + "/live/" + req.StockSymbol); err == nil {
		defer stockResp.Body.Close()
		if stockResp.StatusCode == 200 {
			var stock StockLive
			if json.NewDecoder(stockResp.Body).Decode(&stock) == nil {
				currentPrice = &stock.Price
			}
		}
	}

	// Create new alert
	alertCounter++
	alert := Alert{
		ID:             fmt.Sprintf("alert-%d", alertCounter),
		UserID:         "user-123", // In production, get from authentication
		StockSymbol:    req.StockSymbol,
		StockName:      req.StockName,
		AlertType:      req.AlertType,
		ThresholdPrice: req.ThresholdPrice,
		CurrentPrice:   currentPrice,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	alerts = append(alerts, alert)
	
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, alert)
}

func getAlert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	for _, alert := range alerts {
		if alert.ID == id {
			render.JSON(w, r, alert)
			return
		}
	}
	
	http.Error(w, "Alert not found", http.StatusNotFound)
}

func updateAlert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var req UpdateAlertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i, alert := range alerts {
		if alert.ID == id {
			// Update fields if provided
			if req.AlertType != nil {
				alert.AlertType = *req.AlertType
			}
			if req.ThresholdPrice != nil {
				alert.ThresholdPrice = req.ThresholdPrice
			}
			if req.Status != nil {
				alert.Status = *req.Status
			}
			alert.UpdatedAt = time.Now()
			
			alerts[i] = alert
			render.JSON(w, r, alert)
			return
		}
	}
	
	http.Error(w, "Alert not found", http.StatusNotFound)
}

func deleteAlert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	for i, alert := range alerts {
		if alert.ID == id {
			// Remove alert from slice
			alerts = append(alerts[:i], alerts[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	
	http.Error(w, "Alert not found", http.StatusNotFound)
}

// Background service to check alerts (simplified version)
func startAlertMonitoring() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkAlerts()
		}
	}
}

func checkAlerts() {
	for i, alert := range alerts {
		if alert.Status != "active" || alert.AlertType != "price_threshold" {
			continue
		}

		// Get current stock price
		resp, err := http.Get(GSE_BASE_URL + "/live/" + alert.StockSymbol)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var stock StockLive
		if err := json.NewDecoder(resp.Body).Decode(&stock); err != nil {
			continue
		}

		// Update current price
		alerts[i].CurrentPrice = &stock.Price
		alerts[i].UpdatedAt = time.Now()

		// Check if threshold is met
		if alert.ThresholdPrice != nil && stock.Price >= *alert.ThresholdPrice {
			alerts[i].Status = "triggered"
			fmt.Printf("Alert triggered for %s: Price %.2f reached threshold %.2f\n", 
				alert.StockSymbol, stock.Price, *alert.ThresholdPrice)
			
			// In production, send notification here (email, SMS, push notification, etc.)
		}
	}
}

// Initialize some sample alerts
func initSampleAlerts() {
	alerts = []Alert{
		{
			ID:             "alert-1",
			UserID:         "user-123",
			StockSymbol:    "MTN",
			StockName:      "MTN Ghana",
			AlertType:      "price_threshold",
			ThresholdPrice: func() *float64 { v := 0.85; return &v }(),
			CurrentPrice:   func() *float64 { v := 0.82; return &v }(),
			Status:         "active",
			CreatedAt:      time.Now().Add(-24 * time.Hour),
			UpdatedAt:      time.Now().Add(-1 * time.Hour),
		},
		{
			ID:          "alert-2",
			UserID:      "user-123",
			StockSymbol: "ACCESS",
			StockName:   "Access Bank Ghana Plc",
			AlertType:   "dividend_announcement",
			Status:      "active",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now().Add(-48 * time.Hour),
		},
	}
	alertCounter = 2
}