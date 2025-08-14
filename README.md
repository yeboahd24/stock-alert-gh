# Shares Alert Ghana

A full-stack web application for monitoring Ghana Stock Exchange (GSE) stocks and setting up price alerts.

## Features

- **Real-time Stock Data**: Live stock prices from Ghana Stock Exchange
- **Stock Alerts**: Create price threshold alerts for your favorite stocks
- **Detailed Stock Information**: Company details, financial metrics, and market data
- **Responsive Dashboard**: Modern Material-UI interface
- **Alert Management**: Create, view, edit, and delete stock alerts

## Tech Stack

### Frontend
- React 18 with TypeScript
- Material-UI (MUI) for components
- Vite for build tooling
- Recharts for data visualization

### Backend
- Go with Chi router
- RESTful API design
- Integration with Ghana Stock Exchange API
- Real-time alert monitoring

## API Integration

This application integrates with the Ghana Stock Exchange API:
- Live stock data: `https://dev.kwayisi.org/apis/gse/live`
- Detailed company data: `https://dev.kwayisi.org/apis/gse/equities/{symbol}`

## Local Development

### Prerequisites
- Node.js 18+
- Go 1.21+
- Git

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd shares-alert-ghana
   ```

2. **Backend Setup**
   ```bash
   cd backend
   go mod tidy
   go run .
   ```
   Backend will run on http://localhost:8080

3. **Frontend Setup**
   ```bash
   npm install
   npm run dev
   ```
   Frontend will run on http://localhost:3000

## Docker Deployment

### Using Docker Compose (Recommended)
```bash
docker-compose up --build
```

### Individual Services
```bash
# Backend
cd backend
docker build -t shares-alert-backend .
docker run -p 8080:8080 shares-alert-backend

# Frontend
docker build -t shares-alert-frontend .
docker run -p 80:80 shares-alert-frontend
```

## Render Deployment

This application is configured for deployment on Render using the `render.yaml` file.

### Deploy Steps:

1. **Fork/Clone this repository**

2. **Connect to Render**
   - Go to [Render Dashboard](https://dashboard.render.com)
   - Click "New" → "Blueprint"
   - Connect your GitHub repository
   - Render will automatically detect the `render.yaml` file

3. **Environment Variables**
   The application will automatically configure the API URL for production.

### Manual Deployment (Alternative)

If you prefer manual deployment:

1. **Backend Service**
   - Type: Web Service
   - Environment: Go
   - Build Command: `go build -o main .`
   - Start Command: `./main`
   - Port: 8080

2. **Frontend Service**
   - Type: Static Site
   - Build Command: `npm ci && npm run build`
   - Publish Directory: `dist`
   - Environment Variable: `VITE_API_URL=https://your-backend-url.onrender.com/api/v1`

## API Endpoints

### Stock Endpoints
- `GET /api/v1/stocks` - Get all stocks
- `GET /api/v1/stocks/{symbol}` - Get specific stock
- `GET /api/v1/stocks/{symbol}/details` - Get detailed stock information

### Alert Endpoints
- `GET /api/v1/alerts` - Get all alerts
- `POST /api/v1/alerts` - Create new alert
- `GET /api/v1/alerts/{id}` - Get specific alert
- `PUT /api/v1/alerts/{id}` - Update alert
- `DELETE /api/v1/alerts/{id}` - Delete alert

### Health Check
- `GET /api/v1/health` - Service health status

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support or questions, please open an issue in the GitHub repository.