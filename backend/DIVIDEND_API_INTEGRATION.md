# GSE Dividend API Integration

This document describes the integration with the GSE dividends API that provides real-time dividend yield data for Ghana Stock Exchange stocks.

## Overview

The dividend service has been updated to pull dividend data directly from the GSE dividends API at `https://gse-dividends.onrender.com/stocks`. This provides comprehensive dividend information including yield percentages, stock prices, market caps, and sector classifications.

## New API Endpoints

### 1. Get All GSE Dividend Stocks
**GET** `/api/v1/dividends/gse`

Returns the complete dividend data from the GSE API.

**Response:**
```json
{
  "success": true,
  "data": {
    "timestamp": "2025-08-23T12:35:46.418430055Z",
    "source": "https://simplywall.st/stocks/gh/dividend-yield-high",
    "count": 17,
    "stocks": [
      {
        "symbol": "GCB",
        "name": "GCB Bank",
        "dividend_yield": 2.5,
        "price": "GH₵9.85",
        "market_cap": "GH₵2.5b",
        "country": "Ghana",
        "exchange": "GSE",
        "sector": "Banks",
        "url": "https://simplywall.st/stocks/gh/banks/ghse-gcb/gcb-bank-shares"
      }
    ]
  }
}
```

### 2. Get Dividend Data for Specific Stock
**GET** `/api/v1/dividends/gse/{symbol}`

Returns dividend information for a specific stock symbol.

**Parameters:**
- `symbol` (path): Stock symbol (e.g., "GCB", "GOIL")

**Response:**
```json
{
  "symbol": "GCB",
  "name": "GCB Bank",
  "dividend_yield": 2.5,
  "price": "GH₵9.85",
  "market_cap": "GH₵2.5b",
  "country": "Ghana",
  "exchange": "GSE",
  "sector": "Banks",
  "url": "https://simplywall.st/stocks/gh/banks/ghse-gcb/gcb-bank-shares"
}
```

### 3. Get High Dividend Yield Stocks
**GET** `/api/v1/dividends/high-yield?minYield={threshold}`

Returns stocks with dividend yield above the specified threshold.

**Query Parameters:**
- `minYield` (optional): Minimum dividend yield percentage (default: 0.0)

**Response:**
```json
{
  "minYield": 2.0,
  "count": 3,
  "stocks": [
    {
      "symbol": "GCB",
      "name": "GCB Bank",
      "dividend_yield": 2.5,
      "price": "GH₵9.85",
      "market_cap": "GH₵2.5b",
      "country": "Ghana",
      "exchange": "GSE",
      "sector": "Banks",
      "url": "https://simplywall.st/stocks/gh/banks/ghse-gcb/gcb-bank-shares"
    }
  ]
}
```

## Data Structure

### GSEDividendStock
```go
type GSEDividendStock struct {
    Symbol        string  `json:"symbol"`        // Stock symbol (e.g., "GCB")
    Name          string  `json:"name"`          // Company name
    DividendYield float64 `json:"dividend_yield"` // Dividend yield percentage
    Price         string  `json:"price"`         // Current stock price with currency
    MarketCap     string  `json:"market_cap"`    // Market capitalization with currency
    Country       string  `json:"country"`       // Country (Ghana)
    Exchange      string  `json:"exchange"`      // Exchange (GSE)
    Sector        string  `json:"sector"`        // Business sector
    URL           string  `json:"url"`           // SimplyWall.St URL for more details
}
```

## Service Methods

### DividendService.GetGSEDividendStocks()
Fetches the complete dividend data from the GSE API.

### DividendService.GetDividendStockBySymbol(symbol string)
Retrieves dividend information for a specific stock symbol.

### DividendService.GetHighDividendYieldStocks(minYield float64)
Returns stocks with dividend yield above the specified threshold.

## Error Handling

All endpoints return appropriate HTTP status codes:
- `200 OK`: Successful request
- `400 Bad Request`: Invalid parameters
- `404 Not Found`: Stock symbol not found
- `500 Internal Server Error`: API or server errors

## Usage Examples

### Get all dividend stocks
```bash
curl http://localhost:8080/api/v1/dividends/gse
```

### Get dividend data for GCB Bank
```bash
curl http://localhost:8080/api/v1/dividends/gse/GCB
```

### Get stocks with dividend yield >= 3%
```bash
curl "http://localhost:8080/api/v1/dividends/high-yield?minYield=3.0"
```

## Integration Benefits

1. **Real-time Data**: Direct access to current dividend yields and stock prices
2. **Comprehensive Information**: Includes market cap, sector, and external links
3. **Filtering Capabilities**: Easy filtering by dividend yield thresholds
4. **Reliable Source**: Data sourced from SimplyWall.St via GSE API
5. **Performance**: Optimized HTTP client with connection pooling

## Existing Functionality

The existing dividend announcement system remains unchanged:
- `/api/v1/dividends/` - Get stored dividend announcements
- `/api/v1/dividends/upcoming` - Get upcoming dividend payments
- `POST /api/v1/dividends/` - Create dividend announcements (authenticated)

This integration provides both historical dividend announcements and real-time dividend yield data from the market.