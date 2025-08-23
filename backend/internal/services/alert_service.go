package services

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"shares-alert-backend/internal/models"
	"shares-alert-backend/internal/repository"
)

type AlertService struct {
	alertRepo    *repository.AlertRepository
	userRepo     *repository.UserRepository
	stockService *StockService
	emailService *EmailService
}

func NewAlertService(
	alertRepo *repository.AlertRepository,
	userRepo *repository.UserRepository,
	stockService *StockService,
	emailService *EmailService,
) *AlertService {
	return &AlertService{
		alertRepo:    alertRepo,
		userRepo:     userRepo,
		stockService: stockService,
		emailService: emailService,
	}
}

func (s *AlertService) CreateAlert(userID string, req *models.CreateAlertRequest) (*models.Alert, error) {
	// Validate required fields
	if req.StockSymbol == "" || req.AlertType == "" {
		return nil, fmt.Errorf("stockSymbol and alertType are required")
	}

	// Validate alert type
	validAlertTypes := map[string]bool{
		models.AlertTypePriceThreshold:       true,
		models.AlertTypeIPO:                  true,
		models.AlertTypeDividendAnnouncement: true,
		models.AlertTypeHighDividendYield:    true,
		models.AlertTypeDividendYieldChange:  true,
		models.AlertTypeTargetDividendYield:  true,
	}
	if !validAlertTypes[req.AlertType] {
		return nil, fmt.Errorf("invalid alert type")
	}

	// Validate required fields for different alert types
	if req.AlertType == models.AlertTypePriceThreshold && req.ThresholdPrice == nil {
		return nil, fmt.Errorf("thresholdPrice is required for price_threshold alerts")
	}
	if req.AlertType == models.AlertTypeHighDividendYield && req.ThresholdYield == nil {
		return nil, fmt.Errorf("thresholdYield is required for high_dividend_yield alerts")
	}
	if req.AlertType == models.AlertTypeTargetDividendYield && req.TargetYield == nil {
		return nil, fmt.Errorf("targetYield is required for target_dividend_yield alerts")
	}
	if req.AlertType == models.AlertTypeDividendYieldChange && req.YieldChangeThreshold == nil {
		return nil, fmt.Errorf("yieldChangeThreshold is required for dividend_yield_change alerts")
	}

	// Get current stock price and dividend yield
	var currentPrice *float64
	var currentYield *float64
	if stock, err := s.stockService.GetStock(req.StockSymbol); err == nil {
		currentPrice = &stock.CurrentPrice
	}
	
	// For dividend yield alerts, get current yield from GSE API
	if req.AlertType == models.AlertTypeHighDividendYield || 
	   req.AlertType == models.AlertTypeDividendYieldChange || 
	   req.AlertType == models.AlertTypeTargetDividendYield {
		// We'll need to inject DividendService here or get yield data
		// For now, we'll set it to nil and update it during monitoring
	}

	// Create new alert
	alert := &models.Alert{
		ID:                   uuid.New().String(),
		UserID:               userID,
		StockSymbol:          req.StockSymbol,
		StockName:            req.StockName,
		AlertType:            req.AlertType,
		ThresholdPrice:       req.ThresholdPrice,
		CurrentPrice:         currentPrice,
		ThresholdYield:       req.ThresholdYield,
		CurrentYield:         currentYield,
		TargetYield:          req.TargetYield,
		YieldChangeThreshold: req.YieldChangeThreshold,
		Status:               models.AlertStatusActive,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := s.alertRepo.Create(alert); err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	return alert, nil
}

func (s *AlertService) GetUserAlerts(userID string, filters map[string]interface{}) ([]*models.Alert, error) {
	return s.alertRepo.GetByUserID(userID, filters)
}

func (s *AlertService) GetAlert(alertID, userID string) (*models.Alert, error) {
	alert, err := s.alertRepo.GetByID(alertID)
	if err != nil {
		return nil, err
	}

	// Ensure user owns this alert
	if alert.UserID != userID {
		return nil, fmt.Errorf("alert not found")
	}

	return alert, nil
}

func (s *AlertService) UpdateAlert(alertID, userID string, req *models.UpdateAlertRequest) (*models.Alert, error) {
	alert, err := s.GetAlert(alertID, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.AlertType != nil {
		alert.AlertType = *req.AlertType
	}
	if req.ThresholdPrice != nil {
		alert.ThresholdPrice = req.ThresholdPrice
	}
	if req.ThresholdYield != nil {
		alert.ThresholdYield = req.ThresholdYield
	}
	if req.TargetYield != nil {
		alert.TargetYield = req.TargetYield
	}
	if req.YieldChangeThreshold != nil {
		alert.YieldChangeThreshold = req.YieldChangeThreshold
	}
	if req.Status != nil {
		alert.Status = *req.Status
	}

	if err := s.alertRepo.Update(alert); err != nil {
		return nil, fmt.Errorf("failed to update alert: %w", err)
	}

	// Fetch updated alert
	return s.alertRepo.GetByID(alertID)
}

func (s *AlertService) DeleteAlert(alertID, userID string) error {
	alert, err := s.GetAlert(alertID, userID)
	if err != nil {
		return err
	}

	return s.alertRepo.Delete(alert.ID)
}

func (s *AlertService) GetActiveAlertsByType(alertType string) ([]*models.Alert, error) {
	return s.alertRepo.GetActiveAlertsByType(alertType)
}

func (s *AlertService) StartMonitoring() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	log.Println("Starting alert monitoring service...")

	for {
		select {
		case <-ticker.C:
			if err := s.checkAlerts(); err != nil {
				log.Printf("Error checking alerts: %v", err)
			}
		}
	}
}

func (s *AlertService) checkAlerts() error {
	alerts, err := s.alertRepo.GetActiveAlerts()
	if err != nil {
		return fmt.Errorf("failed to get active alerts: %w", err)
	}

	for _, alert := range alerts {
		if err := s.processAlert(alert); err != nil {
			log.Printf("Error processing alert %s: %v", alert.ID, err)
		}
	}

	return nil
}

func (s *AlertService) processAlert(alert *models.Alert) error {
	// Only process price threshold alerts in this function
	// IPO alerts are handled by IPOService
	if alert.AlertType != models.AlertTypePriceThreshold {
		return nil
	}

	// Get current stock price
	stock, err := s.stockService.GetStock(alert.StockSymbol)
	if err != nil {
		return fmt.Errorf("failed to get stock price for %s: %w", alert.StockSymbol, err)
	}

	// Update current price in alert
	alert.CurrentPrice = &stock.CurrentPrice
	if err := s.alertRepo.Update(alert); err != nil {
		log.Printf("Failed to update current price for alert %s: %v", alert.ID, err)
	}

	// Check if threshold is met
	if alert.ThresholdPrice != nil && stock.CurrentPrice >= *alert.ThresholdPrice {
		return s.triggerAlert(alert, stock.CurrentPrice)
	}

	return nil
}

func (s *AlertService) triggerAlert(alert *models.Alert, currentPrice float64) error {
	// Update alert status to triggered
	if err := s.alertRepo.TriggerAlert(alert.ID); err != nil {
		return fmt.Errorf("failed to trigger alert: %w", err)
	}

	log.Printf("Alert triggered for %s: Price %.2f reached threshold %.2f",
		alert.StockSymbol, currentPrice, *alert.ThresholdPrice)

	// Get user for notification
	user, err := s.userRepo.GetByID(alert.UserID)
	if err != nil {
		log.Printf("Failed to get user for alert notification: %v", err)
		return nil // Don't fail the alert trigger if we can't send notification
	}

	// Check user preferences
	prefs, err := s.userRepo.GetPreferences(user.ID)
	if err != nil {
		log.Printf("Failed to get user preferences, assuming defaults: %v", err)
		// Assume email notifications are enabled by default
	}

	// Send email notification if enabled
	if prefs == nil || prefs.EmailNotifications {
		if err := s.emailService.SendAlertEmail(user, alert); err != nil {
			log.Printf("Failed to send alert email: %v", err)
		}
	}

	return nil
}