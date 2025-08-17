package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"shares-alert-backend/internal/config"
	"shares-alert-backend/internal/models"
	"shares-alert-backend/internal/repository"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	config       *config.AuthConfig
	googleConfig *oauth2.Config
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.AuthConfig) *AuthService {
	fmt.Printf("DEBUG: AuthConfig RedirectURL: %s\n", cfg.RedirectURL)
	googleConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		userRepo:     userRepo,
		config:       cfg,
		googleConfig: googleConfig,
	}
}

func (s *AuthService) GetGoogleAuthURL(state string) string {
	authURL := s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("DEBUG: Google OAuth Config RedirectURL: %s\n", s.googleConfig.RedirectURL)
	fmt.Printf("DEBUG: Generated Auth URL: %s\n", authURL)
	return authURL
}

func (s *AuthService) HandleGoogleCallback(code string) (*models.User, string, error) {
	// Exchange code for token
	token, err := s.googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from Google
	client := s.googleConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, "", fmt.Errorf("failed to decode user info: %w", err)
	}

	// Check if user exists
	user, err := s.userRepo.GetByGoogleID(googleUser.ID)
	if err != nil {
		// User doesn't exist, create new user
		user = &models.User{
			ID:            uuid.New().String(),
			Email:         googleUser.Email,
			Name:          googleUser.Name,
			Picture:       googleUser.Picture,
			GoogleID:      googleUser.ID,
			EmailVerified: googleUser.VerifiedEmail,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, "", fmt.Errorf("failed to create user: %w", err)
		}

		// Create default preferences
		prefs := &models.UserPreferences{
			ID:                    uuid.New().String(),
			UserID:                user.ID,
			EmailNotifications:    true,
			PushNotifications:     true,
			NotificationFrequency: "immediate",
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		if err := s.userRepo.CreatePreferences(prefs); err != nil {
			// Log error but don't fail the login
			fmt.Printf("Failed to create user preferences: %v\n", err)
		}
	} else {
		// Update existing user info
		user.Email = googleUser.Email
		user.Name = googleUser.Name
		user.Picture = googleUser.Picture
		user.EmailVerified = googleUser.VerifiedEmail
		user.UpdatedAt = time.Now()

		if err := s.userRepo.Update(user); err != nil {
			return nil, "", fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Generate JWT token
	jwtToken, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return user, jwtToken, nil
}

func (s *AuthService) GenerateJWT(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.JWTExpirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "shares-alert-backend",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := s.ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	return s.userRepo.GetByID(claims.UserID)
}