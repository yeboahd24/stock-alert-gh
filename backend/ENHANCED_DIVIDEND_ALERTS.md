# Enhanced Dividend Alerts System

This document describes the comprehensive dividend alerts system that combines traditional dividend announcements with real-time dividend yield monitoring using the GSE dividends API.

## Overview

The enhanced dividend alerts system provides multiple types of dividend-related notifications:

1. **Traditional Dividend Announcements** - Manual dividend announcements and payments
2. **High Dividend Yield Alerts** - Notifications when stocks exceed yield thresholds
3. **Target Dividend Yield Alerts** - Notifications when stocks reach specific yield targets
4. **Dividend Yield Change Alerts** - Notifications when yields change significantly

## Alert Types

### 1. Dividend Announcement Alerts
**Type:** `dividend_announcement`

Traditional alerts for manually entered dividend announcements.

**Creation:**
```json
POST /api/v1/alerts
{
  "stockSymbol": "GCB",
  "stockName": "GCB Bank",
  "alertType": "dividend_announcement"
}
```

### 2. High Dividend Yield Alerts
**Type:** `high_dividend_yield`

Triggers when any stock (or specific stock) reaches or exceeds a minimum dividend yield threshold.

**Creation:**
```json
POST /api/v1/alerts
{
  "stockSymbol": "",  // Empty for all stocks, or specific symbol
  "stockName": "All Stocks",
  "alertType": "high_dividend_yield",
  "thresholdYield": 3.0  // Alert when yield >= 3.0%
}
```

**Use Cases:**
- Monitor market for high-yield opportunities
- Track when favorite stocks become attractive for dividend income
- Set different thresholds for different risk profiles

### 3. Target Dividend Yield Alerts
**Type:** `target_dividend_yield`

Triggers when a specific stock reaches a target dividend yield, useful for entry timing.

**Creation:**
```json
POST /api/v1/alerts
{
  "stockSymbol": "GOIL",
  "stockName": "GOIL Company",
  "alertType": "target_dividend_yield",
  "targetYield": 4.5  // Alert when GOIL reaches 4.5% yield
}
```

**Use Cases:**
- Wait for optimal entry points based on yield
- Monitor when stocks reach attractive dividend levels
- Time purchases for maximum dividend income

### 4. Dividend Yield Change Alerts
**Type:** `dividend_yield_change`

Triggers when a stock's dividend yield changes by a significant amount, indicating potential opportunities or risks.

**Creation:**
```json
POST /api/v1/alerts
{
  "stockSymbol": "GCB",
  "stockName": "GCB Bank",
  "alertType": "dividend_yield_change",
  "yieldChangeThreshold": 0.5  // Alert when yield changes by ±0.5%
}
```

**Use Cases:**
- Monitor for sudden yield changes that might indicate company issues
- Detect yield improvements that create opportunities
- Track market volatility affecting dividend yields

## API Endpoints

### Alert Management
- `POST /api/v1/alerts` - Create new dividend yield alert
- `GET /api/v1/alerts` - Get user's alerts (includes yield data)
- `PUT /api/v1/alerts/{id}` - Update alert parameters
- `DELETE /api/v1/alerts/{id}` - Delete alert

### Dividend Data
- `GET /api/v1/dividends/gse` - Get all GSE dividend stocks with yields
- `GET /api/v1/dividends/gse/{symbol}` - Get specific stock dividend data
- `GET /api/v1/dividends/high-yield?minYield=3.0` - Filter high yield stocks

## Monitoring System

### Frequency
- **Dividend Yield Monitoring:** Every 30 minutes
- **Traditional Dividend Monitoring:** Every 6 hours
- **Price Alerts:** Every 30 seconds

### Data Source
Real-time dividend yield data is sourced from the GSE dividends API at `https://gse-dividends.onrender.com/stocks`, which provides:
- Current dividend yields
- Stock prices
- Market capitalizations
- Sector information
- Updated timestamps

### Alert Processing
1. **Data Fetching:** System fetches latest dividend data from GSE API
2. **Alert Matching:** Compares current yields against user alert criteria
3. **Threshold Checking:** Evaluates each alert type's specific conditions
4. **Notification:** Sends email notifications and marks alerts as triggered
5. **State Updates:** Updates current yields and last known yields for tracking

## Email Notifications

### High Dividend Yield Alert Email
- **Subject:** "High Dividend Yield Alert: [Stock Name] ([Symbol])"
- **Content:** Current yield, threshold, and investment opportunity message
- **Design:** Green theme with prominent yield display

### Target Dividend Yield Alert Email
- **Subject:** "Target Dividend Yield Reached: [Stock Name] ([Symbol])"
- **Content:** Current yield, target reached, and timing message
- **Design:** Blue theme with target achievement celebration

### Dividend Yield Change Alert Email
- **Subject:** "Dividend Yield Change Alert: [Stock Name] ([Symbol])"
- **Content:** Current yield, previous yield, change amount and direction
- **Design:** Orange theme with change metrics highlighted

## Database Schema

### New Alert Fields
```sql
ALTER TABLE shares_alert_alerts ADD COLUMN (
  threshold_yield DECIMAL(5,2),        -- For high_dividend_yield alerts
  current_yield DECIMAL(5,2),          -- Current stock dividend yield
  target_yield DECIMAL(5,2),           -- For target_dividend_yield alerts
  yield_change_threshold DECIMAL(5,2), -- For dividend_yield_change alerts
  last_yield DECIMAL(5,2)              -- For tracking yield changes
);
```

### Indexes
- `idx_alerts_dividend_yield` - Performance optimization for yield queries
- `idx_alerts_stock_type` - Optimization for stock and alert type combinations

## Usage Examples

### 1. Monitor High-Yield Opportunities
```bash
# Create alert for any stock with 4%+ yield
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "stockSymbol": "",
    "stockName": "High Yield Opportunities",
    "alertType": "high_dividend_yield",
    "thresholdYield": 4.0
  }'
```

### 2. Wait for Optimal Entry Point
```bash
# Alert when GCB reaches 3% yield
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "stockSymbol": "GCB",
    "stockName": "GCB Bank",
    "alertType": "target_dividend_yield",
    "targetYield": 3.0
  }'
```

### 3. Monitor Yield Volatility
```bash
# Alert on significant yield changes for GOIL
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "stockSymbol": "GOIL",
    "stockName": "GOIL Company",
    "alertType": "dividend_yield_change",
    "yieldChangeThreshold": 0.3
  }'
```

## Integration Benefits

1. **Real-time Data:** Live dividend yields from GSE API
2. **Comprehensive Coverage:** All GSE stocks with dividend information
3. **Multiple Alert Types:** Different strategies for different investment goals
4. **Automated Monitoring:** Continuous background monitoring
5. **Rich Notifications:** Detailed email alerts with context
6. **Performance Optimized:** Efficient database queries and caching
7. **Backward Compatible:** Existing dividend announcement system unchanged

## Best Practices

### For Investors
- **High Yield Alerts:** Set conservative thresholds (2-3%) to avoid noise
- **Target Alerts:** Use for specific stocks you're researching
- **Change Alerts:** Set reasonable thresholds (0.3-0.5%) to catch significant moves
- **Diversification:** Create alerts for different sectors and yield ranges

### For System Administration
- **Monitoring:** Track API response times and error rates
- **Database:** Regular cleanup of triggered alerts
- **Email Limits:** Monitor email sending quotas
- **Performance:** Review query performance with growing alert volumes

## Future Enhancements

1. **WebSocket Integration:** Real-time yield updates in the frontend
2. **Mobile Push Notifications:** Instant alerts on mobile devices
3. **Advanced Analytics:** Yield trend analysis and predictions
4. **Portfolio Integration:** Alerts based on portfolio yield targets
5. **Social Features:** Share high-yield discoveries with other users
6. **API Rate Limiting:** Intelligent caching to reduce external API calls

This enhanced dividend alerts system provides comprehensive dividend yield monitoring capabilities while maintaining the existing dividend announcement functionality, creating a powerful tool for dividend-focused investors on the Ghana Stock Exchange.