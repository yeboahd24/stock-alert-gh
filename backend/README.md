# Shares Alert Ghana - Backend API

A Go backend service for the Ghana Stock Exchange alerts application.

## Features

- **Stock Data**: Proxy to Ghana Stock Exchange API with enhanced data formatting
- **Alerts Management**: Create, read, update, and delete stock price alerts
- **Real-time Monitoring**: Background service to monitor stock prices and trigger alerts
- **CORS Support**: Configured for frontend integration

## API Endpoints

### Stock Endpoints
- `GET /api/v1/stocks` - Get all stocks with live data
- `GET /api/v1/stocks/{symbol}` - Get specific stock live data
- `GET /api/v1/stocks/{symbol}/details` - Get detailed stock information

### Alert Endpoints
- `GET /api/v1/alerts` - Get all alerts (supports filtering by userId and status)
- `POST /api/v1/alerts` - Create a new alert
- `GET /api/v1/alerts/{id}` - Get specific alert
- `PUT /api/v1/alerts/{id}` - Update an alert
- `DELETE /api/v1/alerts/{id}` - Delete an alert

### Health Check
- `GET /api/v1/health` - Service health status

## Running the Server

```bash
cd backend
go mod tidy
go run .
```

The server will start on port 8080.

## Example Usage

### Get all stocks
```bash
curl http://localhost:8080/api/v1/stocks
```

### Get specific stock
```bash
curl http://localhost:8080/api/v1/stocks/ACCESS
```

### Create an alert
```bash
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "stockSymbol": "MTN",
    "stockName": "MTN Ghana",
    "alertType": "price_threshold",
    "thresholdPrice": 0.90
  }'
```

## Data Sources

This API uses the Ghana Stock Exchange API from https://dev.kwayisi.org/apis/gse/ for real-time stock data.