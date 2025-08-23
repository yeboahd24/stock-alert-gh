# Frontend Integration for Enhanced Dividend Alerts

This document provides comprehensive guidance for integrating and testing the enhanced dividend alerts system in the frontend application.

## 🎯 Overview

The frontend integration includes:
- **Enhanced AlertForm** with dividend yield alert types
- **DividendYieldDashboard** for monitoring real-time dividend data
- **DividendYieldAlerts** for managing dividend-based alerts
- **Comprehensive API integration** with the backend
- **Full test coverage** for all components

## 📁 File Structure

```
components/
├── forms/
│   └── AlertForm.tsx                    # Enhanced with dividend yield fields
├── dividend/
│   ├── DividendYieldDashboard.tsx      # Real-time dividend monitoring
│   └── DividendYieldAlerts.tsx         # Alert management interface
└── tables/
    └── AlertsTable.tsx                 # Updated to show yield data

src/services/
└── api.ts                             # Enhanced with dividend API functions

types/
├── enums.ts                           # Updated with new alert types
└── schema.ts                          # Enhanced with dividend yield types

utils/
└── formatters.ts                      # Updated with yield formatting

tests/
└── dividend-alerts.test.tsx           # Comprehensive test suite
```

## 🔧 Component Integration

### 1. Enhanced AlertForm

The `AlertForm` component now supports three new dividend yield alert types:

#### High Dividend Yield Alerts
```tsx
// Usage example
<AlertForm
  open={true}
  onClose={handleClose}
  onSubmit={handleSubmit}
  stocks={availableStocks}
/>

// Form data for high dividend yield alert
{
  stockSymbol: "",  // Empty for all stocks
  stockName: "High Yield Opportunities",
  alertType: AlertType.HIGH_DIVIDEND_YIELD,
  thresholdYield: 3.0
}
```

#### Target Dividend Yield Alerts
```tsx
// Form data for target dividend yield alert
{
  stockSymbol: "GCB",
  stockName: "GCB Bank",
  alertType: AlertType.TARGET_DIVIDEND_YIELD,
  targetYield: 4.0
}
```

#### Dividend Yield Change Alerts
```tsx
// Form data for dividend yield change alert
{
  stockSymbol: "GOIL",
  stockName: "GOIL Company",
  alertType: AlertType.DIVIDEND_YIELD_CHANGE,
  yieldChangeThreshold: 0.5
}
```

### 2. DividendYieldDashboard

Real-time dividend monitoring with GSE API integration:

```tsx
import DividendYieldDashboard from '../components/dividend/DividendYieldDashboard';

// Usage in your main dashboard
<DividendYieldDashboard
  onCreateAlert={(alertData) => {
    // Handle alert creation
    setAlertFormData(alertData);
    setAlertFormOpen(true);
  }}
/>
```

**Features:**
- Real-time dividend yield data from GSE API
- Search and filter functionality
- Summary statistics (total stocks, high yield count, average yield)
- Quick alert creation for high-yield stocks
- Direct links to SimplyWall.St for detailed analysis

### 3. DividendYieldAlerts

Comprehensive alert management interface:

```tsx
import DividendYieldAlerts from '../components/dividend/DividendYieldAlerts';

// Usage in alerts section
<DividendYieldAlerts
  onEditAlert={(alert) => {
    // Handle alert editing
    setEditingAlert(alert);
    setAlertFormOpen(true);
  }}
  onCreateAlert={() => {
    // Handle new alert creation
    setAlertFormOpen(true);
  }}
/>
```

**Features:**
- Display all dividend-related alerts
- Filter by alert type and status
- Pause/resume alerts
- Edit alert parameters
- Delete alerts with confirmation
- Real-time status updates

## 🔌 API Integration

### Dividend API Functions

```typescript
import { dividendApi } from '../src/services/api';

// Get all GSE dividend stocks
const dividendData = await dividendApi.getGSEDividendStocks();

// Get specific stock dividend data
const stockData = await dividendApi.getDividendStockBySymbol('GCB');

// Get high dividend yield stocks
const highYieldStocks = await dividendApi.getHighDividendYieldStocks(3.0);

// Get traditional dividend announcements
const announcements = await dividendApi.getAllDividends();

// Get upcoming dividend payments
const upcoming = await dividendApi.getUpcomingDividends();
```

### Enhanced Alert API

```typescript
import { alertApi } from '../src/services/api';

// Create dividend yield alert
const alert = await alertApi.createAlert({
  stockSymbol: "GCB",
  stockName: "GCB Bank",
  alertType: "high_dividend_yield",
  thresholdYield: 3.0
});

// Update alert with yield parameters
await alertApi.updateAlert(alertId, {
  thresholdYield: 3.5,
  targetYield: 4.0,
  yieldChangeThreshold: 0.5
});
```

## 🎨 UI/UX Features

### Visual Indicators

```tsx
// Yield trend indicators
{stock.dividend_yield >= 3.0 ? (
  <TrendingUp color="success" />
) : stock.dividend_yield > 0 ? (
  <TrendingUp color="primary" />
) : (
  <TrendingDown color="disabled" />
)}

// Alert type icons
const getAlertIcon = (alertType: string) => {
  switch (alertType) {
    case AlertType.HIGH_DIVIDEND_YIELD:
      return <TrendingUp color="success" />;
    case AlertType.TARGET_DIVIDEND_YIELD:
      return <Target color="primary" />;
    case AlertType.DIVIDEND_YIELD_CHANGE:
      return <ChangeCircle color="warning" />;
    default:
      return <NotificationsActive />;
  }
};
```

### Color Coding

- **Green**: High dividend yields (≥3.0%)
- **Blue**: Target dividend yields
- **Orange**: Dividend yield changes
- **Red**: Alert actions (delete, pause)

### Responsive Design

All components are fully responsive with Material-UI breakpoints:

```tsx
<Grid container spacing={3}>
  <Grid item xs={12} sm={6} md={3}>
    {/* Summary cards */}
  </Grid>
</Grid>

<Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
  {/* Responsive filter controls */}
</Stack>
```

## 🧪 Testing Strategy

### Unit Tests

```bash
# Run dividend alerts tests
npm test dividend-alerts.test.tsx

# Run with coverage
npm test -- --coverage dividend-alerts.test.tsx
```

### Test Coverage

The test suite covers:

1. **AlertForm Component**
   - Rendering all dividend yield alert types
   - Form validation for each alert type
   - Successful form submission
   - Error handling

2. **DividendYieldDashboard Component**
   - Data loading and display
   - Search and filter functionality
   - Alert creation actions
   - Error handling and retry

3. **DividendYieldAlerts Component**
   - Alert list display
   - Alert management actions (edit, delete, toggle)
   - Empty state handling
   - Confirmation dialogs

4. **API Integration**
   - All dividend API endpoints
   - Error handling and retries
   - Data formatting and validation

### Integration Tests

```bash
# Run backend integration tests
cd backend
go run test_integration_dividend_alerts.go
```

### E2E Test Scenarios

1. **Create High Dividend Yield Alert**
   ```typescript
   // Test scenario
   1. Navigate to alerts page
   2. Click "Create Alert"
   3. Select "High Dividend Yield" type
   4. Set threshold yield to 3.0%
   5. Submit form
   6. Verify alert appears in list
   ```

2. **Monitor Dividend Dashboard**
   ```typescript
   // Test scenario
   1. Navigate to dividend dashboard
   2. Verify data loads from GSE API
   3. Use search to filter stocks
   4. Set minimum yield filter
   5. Create alert from dashboard
   ```

3. **Manage Existing Alerts**
   ```typescript
   // Test scenario
   1. View dividend alerts list
   2. Pause an active alert
   3. Edit alert parameters
   4. Delete alert with confirmation
   ```

## 🚀 Deployment Integration

### Environment Variables

```bash
# Frontend (.env)
VITE_API_URL=https://your-backend-api.com/api/v1

# Backend
GSE_DIVIDENDS_API_URL=https://gse-dividends.onrender.com/stocks
```

### Build Integration

```json
// package.json scripts
{
  "scripts": {
    "test:dividend": "vitest run dividend-alerts.test.tsx",
    "test:integration": "cd backend && go run test_integration_dividend_alerts.go",
    "build:test": "npm run test:dividend && npm run build"
  }
}
```

## 📊 Performance Considerations

### Data Caching

```typescript
// Cache dividend data for 30 minutes
const DIVIDEND_CACHE_TTL = 30 * 60 * 1000;

// Use React Query for caching
const { data: dividendData, isLoading } = useQuery(
  ['dividends', 'gse'],
  dividendApi.getGSEDividendStocks,
  {
    staleTime: DIVIDEND_CACHE_TTL,
    cacheTime: DIVIDEND_CACHE_TTL,
  }
);
```

### Lazy Loading

```typescript
// Lazy load dividend components
const DividendYieldDashboard = lazy(() => 
  import('../components/dividend/DividendYieldDashboard')
);

const DividendYieldAlerts = lazy(() => 
  import('../components/dividend/DividendYieldAlerts')
);
```

### Pagination

```typescript
// Implement pagination for large datasets
const [page, setPage] = useState(0);
const [rowsPerPage, setRowsPerPage] = useState(25);

const paginatedStocks = filteredStocks.slice(
  page * rowsPerPage,
  page * rowsPerPage + rowsPerPage
);
```

## 🔧 Development Workflow

### 1. Setup Development Environment

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Start backend (in separate terminal)
cd backend
go run cmd/server/main.go
```

### 2. Component Development

```bash
# Create new dividend component
mkdir -p components/dividend
touch components/dividend/NewComponent.tsx

# Add to index
echo "export { default as NewComponent } from './NewComponent';" >> components/dividend/index.ts
```

### 3. Testing Workflow

```bash
# Run tests in watch mode
npm run test:watch

# Run specific test file
npm test dividend-alerts.test.tsx

# Generate coverage report
npm run test:coverage
```

### 4. Integration Testing

```bash
# Test with backend running
npm run test:integration

# Test API endpoints
curl http://localhost:8080/api/v1/dividends/gse
curl http://localhost:8080/api/v1/dividends/high-yield?minYield=3.0
```

## 🐛 Troubleshooting

### Common Issues

1. **API Connection Errors**
   ```typescript
   // Check API base URL configuration
   console.log('API Base URL:', import.meta.env.VITE_API_URL);
   
   // Verify CORS settings in backend
   // Check network tab in browser dev tools
   ```

2. **Type Errors**
   ```typescript
   // Ensure types are properly imported
   import { AlertType } from '../types/enums';
   import { GSEDividendStock } from '../types/schema';
   ```

3. **Test Failures**
   ```bash
   # Clear test cache
   npm run test:clear-cache
   
   # Update snapshots
   npm test -- --updateSnapshot
   ```

### Debug Mode

```typescript
// Enable debug logging
const DEBUG = import.meta.env.DEV;

if (DEBUG) {
  console.log('Dividend data:', dividendData);
  console.log('Alert form data:', formData);
}
```

## 📈 Future Enhancements

### Planned Features

1. **Real-time Updates**
   - WebSocket integration for live dividend yield updates
   - Push notifications for alert triggers

2. **Advanced Analytics**
   - Dividend yield trend charts
   - Sector-wise dividend analysis
   - Historical yield data

3. **Portfolio Integration**
   - Track dividend yields for user portfolios
   - Calculate total dividend income
   - Yield-based portfolio optimization

4. **Mobile App**
   - React Native components
   - Mobile-specific alert management
   - Offline data caching

This comprehensive integration guide ensures smooth development, testing, and deployment of the enhanced dividend alerts system in the frontend application.