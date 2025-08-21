package models

import "time"

type Alert struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"userId" db:"user_id"`
	StockSymbol    string    `json:"stockSymbol" db:"stock_symbol"`
	StockName      string    `json:"stockName" db:"stock_name"`
	AlertType      string    `json:"alertType" db:"alert_type"`
	ThresholdPrice *float64  `json:"thresholdPrice,omitempty" db:"threshold_price"`
	CurrentPrice   *float64  `json:"currentPrice,omitempty" db:"current_price"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
	TriggeredAt    *time.Time `json:"triggeredAt,omitempty" db:"triggered_at"`
}

type CreateAlertRequest struct {
	StockSymbol    string   `json:"stockSymbol"`
	StockName      string   `json:"stockName"`
	AlertType      string   `json:"alertType"`
	ThresholdPrice *float64 `json:"thresholdPrice,omitempty"`
}

type UpdateAlertRequest struct {
	AlertType      *string  `json:"alertType,omitempty"`
	ThresholdPrice *float64 `json:"thresholdPrice,omitempty"`
	Status         *string  `json:"status,omitempty"`
}

type IPOAnnouncement struct {
	ID          string    `json:"id" db:"id"`
	CompanyName string    `json:"companyName" db:"company_name"`
	Symbol      string    `json:"symbol" db:"symbol"`
	Sector      string    `json:"sector" db:"sector"`
	OfferPrice  float64   `json:"offerPrice" db:"offer_price"`
	ListingDate time.Time `json:"listingDate" db:"listing_date"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateIPORequest struct {
	CompanyName string    `json:"companyName"`
	Symbol      string    `json:"symbol"`
	Sector      string    `json:"sector"`
	OfferPrice  float64   `json:"offerPrice"`
	ListingDate time.Time `json:"listingDate"`
}

// Alert types
const (
	AlertTypePriceThreshold       = "price_threshold"
	AlertTypeIPO                  = "ipo_alert"
	AlertTypeDividendAnnouncement = "dividend_announcement"
)

// Alert statuses
const (
	AlertStatusActive    = "active"
	AlertStatusTriggered = "triggered"
	AlertStatusPaused    = "paused"
	AlertStatusDeleted   = "deleted"
)

// IPO statuses
const (
	IPOStatusAnnounced = "announced"
	IPOStatusListed    = "listed"
	IPOStatusCancelled = "cancelled"
)