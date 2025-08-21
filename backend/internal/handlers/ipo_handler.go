package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"

	"shares-alert-backend/internal/models"
	"shares-alert-backend/internal/services"
)

type IPOHandler struct {
	ipoService *services.IPOService
}

func NewIPOHandler(ipoService *services.IPOService) *IPOHandler {
	return &IPOHandler{
		ipoService: ipoService,
	}
}

func (h *IPOHandler) CreateIPO(w http.ResponseWriter, r *http.Request) {
	var req models.CreateIPORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ipo, err := h.ipoService.CreateIPOAnnouncement(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, ipo)
}

func (h *IPOHandler) GetAllIPOs(w http.ResponseWriter, r *http.Request) {
	ipos, err := h.ipoService.GetAllIPOs()
	if err != nil {
		http.Error(w, "Failed to fetch IPOs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, ipos)
}

func (h *IPOHandler) GetUpcomingIPOs(w http.ResponseWriter, r *http.Request) {
	ipos, err := h.ipoService.GetUpcomingIPOs()
	if err != nil {
		http.Error(w, "Failed to fetch upcoming IPOs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, ipos)
}