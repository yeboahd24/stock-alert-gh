package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"shares-alert-backend/internal/models"
	"shares-alert-backend/internal/services"
)

type DividendHandler struct {
	dividendService *services.DividendService
}

func NewDividendHandler(dividendService *services.DividendService) *DividendHandler {
	return &DividendHandler{
		dividendService: dividendService,
	}
}

func (h *DividendHandler) CreateDividend(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDividendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dividend, err := h.dividendService.CreateDividendAnnouncement(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, dividend)
}

func (h *DividendHandler) GetAllDividends(w http.ResponseWriter, r *http.Request) {
	dividends, err := h.dividendService.GetAllDividends()
	if err != nil {
		http.Error(w, "Failed to fetch dividends: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dividends)
}

func (h *DividendHandler) GetUpcomingDividends(w http.ResponseWriter, r *http.Request) {
	dividends, err := h.dividendService.GetUpcomingDividends()
	if err != nil {
		http.Error(w, "Failed to fetch upcoming dividends: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dividends)
}

// GetGSEDividendStocks fetches dividend data from the GSE API
func (h *DividendHandler) GetGSEDividendStocks(w http.ResponseWriter, r *http.Request) {
	dividendData, err := h.dividendService.GetGSEDividendStocks()
	if err != nil {
		http.Error(w, "Failed to fetch GSE dividend data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dividendData)
}

// GetDividendStockBySymbol fetches dividend information for a specific stock
func (h *DividendHandler) GetDividendStockBySymbol(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		http.Error(w, "Stock symbol is required", http.StatusBadRequest)
		return
	}

	stock, err := h.dividendService.GetDividendStockBySymbol(symbol)
	if err != nil {
		http.Error(w, "Failed to fetch dividend data for symbol: "+err.Error(), http.StatusNotFound)
		return
	}

	render.JSON(w, r, stock)
}

// GetHighDividendYieldStocks returns stocks with high dividend yields
func (h *DividendHandler) GetHighDividendYieldStocks(w http.ResponseWriter, r *http.Request) {
	minYieldStr := r.URL.Query().Get("minYield")
	minYield := 0.0 // Default minimum yield
	
	if minYieldStr != "" {
		var err error
		minYield, err = strconv.ParseFloat(minYieldStr, 64)
		if err != nil {
			http.Error(w, "Invalid minYield parameter", http.StatusBadRequest)
			return
		}
	}

	stocks, err := h.dividendService.GetHighDividendYieldStocks(minYield)
	if err != nil {
		http.Error(w, "Failed to fetch high dividend yield stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"minYield": minYield,
		"count":    len(stocks),
		"stocks":   stocks,
	})
}