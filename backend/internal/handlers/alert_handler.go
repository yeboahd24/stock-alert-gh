package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"shares-alert-backend/internal/models"
	"shares-alert-backend/internal/services"
)

type AlertHandler struct {
	alertService *services.AlertService
}

func NewAlertHandler(alertService *services.AlertService) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
	}
}

func (h *AlertHandler) GetAlerts(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse query parameters for filtering
	filters := make(map[string]interface{})
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}
	if stockSymbol := r.URL.Query().Get("stockSymbol"); stockSymbol != "" {
		filters["stock_symbol"] = stockSymbol
	}
	if alertType := r.URL.Query().Get("alertType"); alertType != "" {
		filters["alert_type"] = alertType
	}

	alerts, err := h.alertService.GetUserAlerts(user.ID, filters)
	if err != nil {
		http.Error(w, "Failed to fetch alerts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, alerts)
}

func (h *AlertHandler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req models.CreateAlertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	alert, err := h.alertService.CreateAlert(user.ID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, alert)
}

func (h *AlertHandler) GetAlert(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	alertID := chi.URLParam(r, "id")
	if alertID == "" {
		http.Error(w, "Alert ID is required", http.StatusBadRequest)
		return
	}

	alert, err := h.alertService.GetAlert(alertID, user.ID)
	if err != nil {
		http.Error(w, "Alert not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, alert)
}

func (h *AlertHandler) UpdateAlert(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	alertID := chi.URLParam(r, "id")
	if alertID == "" {
		http.Error(w, "Alert ID is required", http.StatusBadRequest)
		return
	}

	var req models.UpdateAlertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	alert, err := h.alertService.UpdateAlert(alertID, user.ID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	render.JSON(w, r, alert)
}

func (h *AlertHandler) DeleteAlert(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	alertID := chi.URLParam(r, "id")
	if alertID == "" {
		http.Error(w, "Alert ID is required", http.StatusBadRequest)
		return
	}

	if err := h.alertService.DeleteAlert(alertID, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}