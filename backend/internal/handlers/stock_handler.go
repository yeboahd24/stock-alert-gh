package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"shares-alert-backend/internal/services"
)

type StockHandler struct {
	stockService *services.StockService
}

func NewStockHandler(stockService *services.StockService) *StockHandler {
	return &StockHandler{
		stockService: stockService,
	}
}

func (h *StockHandler) GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := h.stockService.GetAllStocks()
	if err != nil {
		http.Error(w, "Failed to fetch stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, stocks)
}

func (h *StockHandler) GetStock(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		http.Error(w, "Stock symbol is required", http.StatusBadRequest)
		return
	}

	stock, err := h.stockService.GetStock(symbol)
	if err != nil {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, stock)
}

func (h *StockHandler) GetStockDetails(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		http.Error(w, "Stock symbol is required", http.StatusBadRequest)
		return
	}

	stock, err := h.stockService.GetStockDetails(symbol)
	if err != nil {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, stock)
}