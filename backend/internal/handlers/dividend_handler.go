package handlers

import (
	"encoding/json"
	"net/http"

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