package models

import "time"

type Alert struct {
	ID                   string    `json:"id" db:"id"`
	UserID               string    `json:"userId" db:"user_id"`
	StockSymbol          string    `json:"stockSymbol" db:"stock_symbol"`
	StockName            string    `json:"stockName" db:"stock_name"`
	AlertType            string    `json:"alertType" db:"alert_type"`
	ThresholdPrice       *float64  `json:"thresholdPrice,omitempty" db:"threshold_price"`
	CurrentPrice         *float64  `json:"currentPrice,omitempty" db:"current_price"`
	ThresholdYield       *float64  `json:"thresholdYield,omitempty" db:"threshold_yield"`
	CurrentYield         *float64  `json:"currentYield,omitempty" db:"current_yield"`
	TargetYield          *float64  `json:"targetYield,omitempty" db:"target_yield"`
	YieldChangeThreshold *float64  `json:"yieldChangeThreshold,omitempty" db:"yield_change_threshold"`
	LastYield            *float64  `json:"lastYield,omitempty" db:"last_yield"`
	Status               string    `json:"status" db:"status"`
	CreatedAt            time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time `json:"updatedAt" db:"updated_at"`
	TriggeredAt          *time.Time `json:"triggeredAt,omitempty" db:"triggered_at"`
}

type CreateAlertRequest struct {
	StockSymbol          string   `json:"stockSymbol"`
	StockName            string   `json:"stockName"`
	AlertType            string   `json:"alertType"`
	ThresholdPrice       *float64 `json:"thresholdPrice,omitempty"`
	ThresholdYield       *float64 `json:"thresholdYield,omitempty"`
	TargetYield          *float64 `json:"targetYield,omitempty"`
	YieldChangeThreshold *float64 `json:"yieldChangeThreshold,omitempty"`
}

type UpdateAlertRequest struct {
	AlertType            *string  `json:"alertType,omitempty"`
	ThresholdPrice       *float64 `json:"thresholdPrice,omitempty"`
	ThresholdYield       *float64 `json:"thresholdYield,omitempty"`
	TargetYield          *float64 `json:"targetYield,omitempty"`
	YieldChangeThreshold *float64 `json:"yieldChangeThreshold,omitempty"`
	Status               *string  `json:"status,omitempty"`
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
	AlertTypeHighDividendYield    = "high_dividend_yield"
	AlertTypeDividendYieldChange  = "dividend_yield_change"
	AlertTypeTargetDividendYield  = "target_dividend_yield"
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

// Dividend statuses
const (
	DividendStatusAnnounced = "announced"
	DividendStatusPaid      = "paid"
	DividendStatusCancelled = "cancelled"
)

type DividendAnnouncement struct {
	ID           string    `json:"id" db:"id"`
	StockSymbol  string    `json:"stockSymbol" db:"stock_symbol"`
	StockName    string    `json:"stockName" db:"stock_name"`
	DividendType string    `json:"dividendType" db:"dividend_type"`
	Amount       float64   `json:"amount" db:"amount"`
	Currency     string    `json:"currency" db:"currency"`
	ExDate       time.Time `json:"exDate" db:"ex_date"`
	PaymentDate  time.Time `json:"paymentDate" db:"payment_date"`
	Status       string    `json:"status" db:"status"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateDividendRequest struct {
	StockSymbol  string    `json:"stockSymbol"`
	StockName    string    `json:"stockName"`
	DividendType string    `json:"dividendType"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	ExDate       time.Time `json:"exDate"`
	PaymentDate  time.Time `json:"paymentDate"`
}