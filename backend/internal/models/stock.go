package models

import "time"

// Stock data structures based on the Ghana Stock Exchange API
type StockLive struct {
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Change float64 `json:"change"`
	Volume int64   `json:"volume"`
}

type Company struct {
	Address   string   `json:"address"`
	Directors []string `json:"directors"`
	Email     string   `json:"email"`
	Facsimile *string  `json:"facsimile"`
	Industry  string   `json:"industry"`
	Name      string   `json:"name"`
	Sector    string   `json:"sector"`
	Telephone string   `json:"telephone"`
	Website   string   `json:"website"`
}

type StockEquity struct {
	Capital float64  `json:"capital"`
	Company Company  `json:"company"`
	DPS     *float64 `json:"dps"`
	EPS     *float64 `json:"eps"`
	Name    string   `json:"name"`
	Price   float64  `json:"price"`
	Shares  int64    `json:"shares"`
}

// Enhanced stock data for our application
type EnhancedStock struct {
	Symbol           string    `json:"symbol"`
	Name             string    `json:"name"`
	CurrentPrice     float64   `json:"currentPrice"`
	PreviousClose    float64   `json:"previousClose"`
	Change           float64   `json:"change"`
	ChangePercent    float64   `json:"changePercent"`
	Volume           int64     `json:"volume"`
	LastUpdated      time.Time `json:"lastUpdated"`
	MarketCap        *float64  `json:"marketCap,omitempty"`
	Sector           *string   `json:"sector,omitempty"`
	Industry         *string   `json:"industry,omitempty"`
}

// Detailed stock data including company information
type DetailedStock struct {
	Symbol           string    `json:"symbol"`
	Name             string    `json:"name"`
	CurrentPrice     float64   `json:"currentPrice"`
	PreviousClose    float64   `json:"previousClose"`
	Change           float64   `json:"change"`
	ChangePercent    float64   `json:"changePercent"`
	Volume           int64     `json:"volume"`
	LastUpdated      time.Time `json:"lastUpdated"`
	MarketCap        float64   `json:"marketCap"`
	Shares           int64     `json:"shares"`
	Sector           string    `json:"sector"`
	Industry         string    `json:"industry"`
	DPS              *float64  `json:"dps"`
	EPS              *float64  `json:"eps"`
	Company          Company   `json:"company"`
}

// GSE Dividend API response structures
type GSEDividendStock struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	DividendYield float64 `json:"dividend_yield"`
	Price         string  `json:"price"`
	MarketCap     string  `json:"market_cap"`
	Country       string  `json:"country"`
	Exchange      string  `json:"exchange"`
	Sector        string  `json:"sector"`
	URL           string  `json:"url"`
}

type GSEDividendData struct {
	Timestamp string             `json:"timestamp"`
	Source    string             `json:"source"`
	Count     int                `json:"count"`
	Stocks    []GSEDividendStock `json:"stocks"`
}

type GSEDividendResponse struct {
	Success bool            `json:"success"`
	Data    GSEDividendData `json:"data"`
}