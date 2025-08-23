package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"shares-alert-backend/internal/httpclient"
	"shares-alert-backend/internal/models"
	"shares-alert-backend/internal/repository"
)

type DividendService struct {
	dividendRepo *repository.DividendRepository
	alertRepo    *repository.AlertRepository
	userRepo     *repository.UserRepository
	emailService *EmailService
}

func NewDividendService(
	dividendRepo *repository.DividendRepository,
	alertRepo *repository.AlertRepository,
	userRepo *repository.UserRepository,
	emailService *EmailService,
) *DividendService {
	return &DividendService{
		dividendRepo: dividendRepo,
		alertRepo:    alertRepo,
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *DividendService) CreateDividendAnnouncement(req *models.CreateDividendRequest) (*models.DividendAnnouncement, error) {
	dividend := &models.DividendAnnouncement{
		ID:           uuid.New().String(),
		StockSymbol:  req.StockSymbol,
		StockName:    req.StockName,
		DividendType: req.DividendType,
		Amount:       req.Amount,
		Currency:     req.Currency,
		ExDate:       req.ExDate,
		PaymentDate:  req.PaymentDate,
		Status:       models.DividendStatusAnnounced,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.dividendRepo.Create(dividend); err != nil {
		return nil, fmt.Errorf("failed to create dividend announcement: %w", err)
	}

	// Trigger dividend alerts
	go s.triggerDividendAlerts(dividend)

	return dividend, nil
}

func (s *DividendService) GetAllDividends() ([]*models.DividendAnnouncement, error) {
	return s.dividendRepo.GetAll()
}

func (s *DividendService) GetUpcomingDividends() ([]*models.DividendAnnouncement, error) {
	return s.dividendRepo.GetUpcoming()
}

// GetGSEDividendStocks fetches dividend data directly from the GSE dividends API
func (s *DividendService) GetGSEDividendStocks() (*models.GSEDividendResponse, error) {
	client := httpclient.GetDefaultClient()
	
	resp, err := client.Get("https://gse-dividends.onrender.com/stocks")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GSE dividend data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GSE API returned status %d", resp.StatusCode)
	}

	var dividendResponse models.GSEDividendResponse
	if err := json.NewDecoder(resp.Body).Decode(&dividendResponse); err != nil {
		return nil, fmt.Errorf("failed to decode GSE dividend response: %w", err)
	}

	if !dividendResponse.Success {
		return nil, fmt.Errorf("GSE API returned unsuccessful response")
	}

	return &dividendResponse, nil
}

// GetDividendStockBySymbol fetches dividend information for a specific stock symbol
func (s *DividendService) GetDividendStockBySymbol(symbol string) (*models.GSEDividendStock, error) {
	dividendData, err := s.GetGSEDividendStocks()
	if err != nil {
		return nil, err
	}

	for _, stock := range dividendData.Data.Stocks {
		if stock.Symbol == symbol {
			return &stock, nil
		}
	}

	return nil, fmt.Errorf("stock with symbol %s not found in dividend data", symbol)
}

// GetHighDividendYieldStocks returns stocks with dividend yield above a threshold
func (s *DividendService) GetHighDividendYieldStocks(minYield float64) ([]models.GSEDividendStock, error) {
	dividendData, err := s.GetGSEDividendStocks()
	if err != nil {
		return nil, err
	}

	var highYieldStocks []models.GSEDividendStock
	for _, stock := range dividendData.Data.Stocks {
		if stock.DividendYield >= minYield {
			highYieldStocks = append(highYieldStocks, stock)
		}
	}

	return highYieldStocks, nil
}

func (s *DividendService) StartDividendMonitoring() {
	ticker := time.NewTicker(6 * time.Hour) // Check every 6 hours
	defer ticker.Stop()

	log.Println("Starting dividend monitoring service...")

	for {
		select {
		case <-ticker.C:
			if err := s.checkDividendPayments(); err != nil {
				log.Printf("Error checking dividend payments: %v", err)
			}
		}
	}
}

func (s *DividendService) checkDividendPayments() error {
	upcomingDividends, err := s.dividendRepo.GetUpcoming()
	if err != nil {
		return fmt.Errorf("failed to get upcoming dividends: %w", err)
	}

	now := time.Now()
	for _, dividend := range upcomingDividends {
		// Check if payment date has passed
		if dividend.PaymentDate.Before(now) && dividend.Status == models.DividendStatusAnnounced {
			// Update status to paid
			if err := s.dividendRepo.UpdateStatus(dividend.ID, models.DividendStatusPaid); err != nil {
				log.Printf("Failed to update dividend status for %s: %v", dividend.StockSymbol, err)
				continue
			}

			// Trigger payment alerts
			go s.triggerDividendPaymentAlerts(dividend)
		}
	}

	return nil
}

func (s *DividendService) triggerDividendAlerts(dividend *models.DividendAnnouncement) {
	// Get all active dividend alerts
	alerts, err := s.alertRepo.GetActiveAlertsByType(models.AlertTypeDividendAnnouncement)
	if err != nil {
		log.Printf("Failed to get dividend alerts: %v", err)
		return
	}

	for _, alert := range alerts {
		// For stock-specific alerts, check symbol match
		if alert.StockSymbol != "" && alert.StockSymbol != dividend.StockSymbol {
			continue
		}

		if err := s.notifyDividendAlert(alert, dividend, "announced"); err != nil {
			log.Printf("Failed to notify dividend alert %s: %v", alert.ID, err)
		}
	}
}

func (s *DividendService) triggerDividendPaymentAlerts(dividend *models.DividendAnnouncement) {
	// Get all active dividend alerts for this stock
	alerts, err := s.alertRepo.GetActiveAlertsByType(models.AlertTypeDividendAnnouncement)
	if err != nil {
		log.Printf("Failed to get dividend alerts: %v", err)
		return
	}

	for _, alert := range alerts {
		if alert.StockSymbol != "" && alert.StockSymbol != dividend.StockSymbol {
			continue
		}

		if err := s.notifyDividendAlert(alert, dividend, "paid"); err != nil {
			log.Printf("Failed to notify dividend payment alert %s: %v", alert.ID, err)
		}
	}
}

func (s *DividendService) notifyDividendAlert(alert *models.Alert, dividend *models.DividendAnnouncement, eventType string) error {
	// Get user for notification
	user, err := s.userRepo.GetByID(alert.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check user preferences
	prefs, err := s.userRepo.GetPreferences(user.ID)
	if err != nil {
		log.Printf("Failed to get user preferences, assuming defaults: %v", err)
	}

	// Send email notification if enabled
	if prefs == nil || prefs.EmailNotifications {
		if err := s.emailService.SendDividendAlertEmail(user, alert, dividend, eventType); err != nil {
			return fmt.Errorf("failed to send dividend alert email: %w", err)
		}
	}

	// Mark alert as triggered for announcements
	if eventType == "announced" {
		if err := s.alertRepo.TriggerAlert(alert.ID); err != nil {
			log.Printf("Failed to mark dividend alert as triggered: %v", err)
		}
	}

	return nil
}