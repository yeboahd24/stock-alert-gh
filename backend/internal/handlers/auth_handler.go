package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"shares-alert-backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

type AuthResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

type GoogleAuthRequest struct {
	Code string `json:"code"`
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) GetGoogleAuthURL(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		state = "random-state-string" // In production, generate a proper state
	}

	authURL := h.authService.GetGoogleAuthURL(state)
	render.JSON(w, r, map[string]string{"authUrl": authURL})
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	var req GoogleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Code == "" {
		http.Error(w, "Authorization code is required", http.StatusBadRequest)
		return
	}

	user, token, err := h.authService.HandleGoogleCallback(req.Code)
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	response := AuthResponse{
		User:  user,
		Token: token,
	}

	render.JSON(w, r, response)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// For JWT-based auth, logout is handled client-side by removing the token
	// Here we just return a success response
	render.JSON(w, r, map[string]string{"message": "Logged out successfully"})
}

// Middleware to authenticate requests
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		user, err := h.authService.GetUserFromToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user to request context
		ctx := r.Context()
		ctx = setUserInContext(ctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Optional middleware - doesn't fail if no auth provided
func (h *AuthHandler) OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				if user, err := h.authService.GetUserFromToken(token); err == nil {
					ctx := r.Context()
					ctx = setUserInContext(ctx, user)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}