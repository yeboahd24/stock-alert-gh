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