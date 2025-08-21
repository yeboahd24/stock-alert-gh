package services

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"shares-alert-backend/internal/models"
	"shares-alert-backend/internal/repository"
)

type IPOService struct {
	ipoRepo      *repository.IPORepository
	alertRepo    *repository.AlertRepository
	userRepo     *repository.UserRepository
	emailService *EmailService
}

func NewIPOService(
	ipoRepo *repository.IPORepository,
	alertRepo *repository.AlertRepository,
	userRepo *repository.UserRepository,
	emailService *EmailService,
) *IPOService {
	return &IPOService{
		ipoRepo:      ipoRepo,
		alertRepo:    alertRepo,
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *IPOService) CreateIPOAnnouncement(req *models.CreateIPORequest) (*models.IPOAnnouncement, error) {
	ipo := &models.IPOAnnouncement{
		ID:          uuid.New().String(),
		CompanyName: req.CompanyName,
		Symbol:      req.Symbol,
		Sector:      req.Sector,
		OfferPrice:  req.OfferPrice,
		ListingDate: req.ListingDate,
		Status:      models.IPOStatusAnnounced,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.ipoRepo.Create(ipo); err != nil {
		return nil, fmt.Errorf("failed to create IPO announcement: %w", err)
	}

	// Trigger IPO alerts
	go s.triggerIPOAlerts(ipo)

	return ipo, nil
}

func (s *IPOService) GetAllIPOs() ([]*models.IPOAnnouncement, error) {
	return s.ipoRepo.GetAll()
}

func (s *IPOService) GetUpcomingIPOs() ([]*models.IPOAnnouncement, error) {
	return s.ipoRepo.GetUpcoming()
}

func (s *IPOService) StartIPOMonitoring() {
	ticker := time.NewTicker(1 * time.Hour) // Check every hour
	defer ticker.Stop()

	log.Println("Starting IPO monitoring service...")

	for {
		select {
		case <-ticker.C:
			if err := s.checkIPOListings(); err != nil {
				log.Printf("Error checking IPO listings: %v", err)
			}
		}
	}
}

func (s *IPOService) checkIPOListings() error {
	upcomingIPOs, err := s.ipoRepo.GetUpcoming()
	if err != nil {
		return fmt.Errorf("failed to get upcoming IPOs: %w", err)
	}

	now := time.Now()
	for _, ipo := range upcomingIPOs {
		// Check if listing date has passed
		if ipo.ListingDate.Before(now) && ipo.Status == models.IPOStatusAnnounced {
			// Update status to listed
			if err := s.ipoRepo.UpdateStatus(ipo.ID, models.IPOStatusListed); err != nil {
				log.Printf("Failed to update IPO status for %s: %v", ipo.Symbol, err)
				continue
			}

			// Trigger listing alerts
			go s.triggerIPOListingAlerts(ipo)
		}
	}

	return nil
}

func (s *IPOService) triggerIPOAlerts(ipo *models.IPOAnnouncement) {
	// Get all active IPO alerts
	alerts, err := s.alertRepo.GetActiveAlertsByType(models.AlertTypeIPO)
	if err != nil {
		log.Printf("Failed to get IPO alerts: %v", err)
		return
	}

	for _, alert := range alerts {
		// For general IPO alerts, notify all users
		if err := s.notifyIPOAlert(alert, ipo, "announced"); err != nil {
			log.Printf("Failed to notify IPO alert %s: %v", alert.ID, err)
		}
	}
}

func (s *IPOService) triggerIPOListingAlerts(ipo *models.IPOAnnouncement) {
	// Get all active IPO alerts
	alerts, err := s.alertRepo.GetActiveAlertsByType(models.AlertTypeIPO)
	if err != nil {
		log.Printf("Failed to get IPO alerts: %v", err)
		return
	}

	for _, alert := range alerts {
		// Notify about listing
		if err := s.notifyIPOAlert(alert, ipo, "listed"); err != nil {
			log.Printf("Failed to notify IPO listing alert %s: %v", alert.ID, err)
		}
	}
}

func (s *IPOService) notifyIPOAlert(alert *models.Alert, ipo *models.IPOAnnouncement, eventType string) error {
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
		if err := s.emailService.SendIPOAlertEmail(user, alert, ipo, eventType); err != nil {
			return fmt.Errorf("failed to send IPO alert email: %w", err)
		}
	}

	// Mark alert as triggered for announcements
	if eventType == "announced" {
		if err := s.alertRepo.TriggerAlert(alert.ID); err != nil {
			log.Printf("Failed to mark IPO alert as triggered: %v", err)
		}
	}

	return nil
}