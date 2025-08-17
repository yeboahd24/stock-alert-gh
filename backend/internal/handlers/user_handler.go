package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"

	"shares-alert-backend/internal/models"
	"shares-alert-backend/internal/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (h *UserHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	prefs, err := h.userRepo.GetPreferences(user.ID)
	if err != nil {
		// Return default preferences if none exist
		defaultPrefs := &models.UserPreferences{
			UserID:                user.ID,
			EmailNotifications:    true,
			PushNotifications:     true,
			NotificationFrequency: "immediate",
		}
		render.JSON(w, r, defaultPrefs)
		return
	}

	render.JSON(w, r, prefs)
}

func (h *UserHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req models.UserPreferences
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the user ID matches
	req.UserID = user.ID

	// Try to update existing preferences
	if err := h.userRepo.UpdatePreferences(&req); err != nil {
		// If update fails, try to create new preferences
		req.ID = user.ID + "-prefs" // Simple ID generation
		if err := h.userRepo.CreatePreferences(&req); err != nil {
			http.Error(w, "Failed to save preferences: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Fetch updated preferences
	prefs, err := h.userRepo.GetPreferences(user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch updated preferences", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, prefs)
}