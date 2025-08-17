# Shares Alert Ghana - Backend API v2.0

A refactored Go backend service for the Ghana Stock Exchange alerts application with improved architecture, Google OAuth authentication, email notifications, and SQLite database support.

## Architecture

The backend has been completely refactored with a clean, modular architecture:

```
backend/
├── cmd/server/           # Application entry point
├── internal/
│   ├── app/             # Application setup and initialization
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and migrations
│   ├── handlers/        # HTTP handlers (controllers)
│   ├── models/          # Data models and structures
│   ├── repository/      # Data access layer
│   └── services/        # Business logic layer
├── data/                # SQLite database files (auto-created)
└── migrations/          # Database migration files
```

## Features

### Core Features
- **Stock Data**: Real-time Ghana Stock Exchange data with fallback to mock data
- **User Authentication**: Google OAuth 2.0 integration
- **Alerts Management**: Create, read, update, and delete stock price alerts
- **Email Notifications**: Automated email alerts when thresholds are met
- **User Preferences**: Customizable notification settings

### Technical Features
- **Database Abstraction**: Easy migration from SQLite to PostgreSQL/MySQL
- **JWT Authentication**: Secure token-based authentication
- **Background Monitoring**: Automated alert checking service
- **CORS Support**: Configured for frontend integration
- **Graceful Error Handling**: Comprehensive error handling and logging
- **Environment Configuration**: Flexible configuration management

## Prerequisites

- Go 1.21 or higher
- Google OAuth 2.0 credentials (for authentication)
- Gmail account with app password (for email notifications)

## Setup

### 1. Clone and Install Dependencies

```bash
cd backend
go mod tidy
```

### 2. Environment Configuration

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Google OAuth (Required)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# Email Configuration (Optional but recommended)
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_gmail_app_password
FROM_EMAIL=your_email@gmail.com

# JWT Secret (Change in production)
JWT_SECRET=your-very-secure-jwt-secret
```

### 3. Google OAuth Setup

1. Go to Google Cloud Console
2. Create a new project or select existing one
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URIs:
   - `http://localhost:3000/auth/callback` (development)
   - `https://yourdomain.com/auth/callback` (production)

### 4. Gmail App Password Setup

1. Enable 2-Factor Authentication on your Gmail account
2. Generate an App Password:
   - Go to Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate password for "Mail"
   - Use this password in `SMTP_PASSWORD`

## Running the Server

### Development
```bash
go run cmd/server/main.go
```

### Production Build
```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

The server will start on port 10000 (or PORT environment variable).

## API Documentation

### Authentication Endpoints

#### Get Google Auth URL
```http
GET /api/v1/auth/google?state=random-string
```

#### Google OAuth Callback
```http
POST /api/v1/auth/google/callback
Content-Type: application/json

{
  "code": "google_auth_code"
}
```

#### Get User Profile
```http
GET /api/v1/auth/profile
Authorization: Bearer <jwt_token>
```

### Stock Endpoints

#### Get All Stocks
```http
GET /api/v1/stocks
```

#### Get Specific Stock
```http
GET /api/v1/stocks/{symbol}
```

#### Get Stock Details
```http
GET /api/v1/stocks/{symbol}/details
```

### Alert Endpoints (Authenticated)

#### Get User Alerts
```http
GET /api/v1/alerts?status=active&stockSymbol=MTN
Authorization: Bearer <jwt_token>
```

#### Create Alert
```http
POST /api/v1/alerts
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "stockSymbol": "MTN",
  "stockName": "MTN Ghana",
  "alertType": "price_threshold",
  "thresholdPrice": 0.90
}
```

#### Update Alert
```http
PUT /api/v1/alerts/{id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "status": "paused",
  "thresholdPrice": 0.95
}
```

#### Delete Alert
```http
DELETE /api/v1/alerts/{id}
Authorization: Bearer <jwt_token>
```

### User Preferences

#### Get Preferences
```http
GET /api/v1/user/preferences
Authorization: Bearer <jwt_token>
```

#### Update Preferences
```http
PUT /api/v1/user/preferences
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "emailNotifications": true,
  "pushNotifications": false,
  "notificationFrequency": "immediate"
}
```

## Database

### SQLite (Default)
The application uses SQLite by default, which is perfect for development and small to medium deployments. The database file is created automatically at `./data/shares_alert.db`.

### Migrating to PostgreSQL

When you're ready to scale, simply update your `.env`:

```env
DB_TYPE=postgres
DB_HOST=your_postgres_host
DB_PORT=5432
DB_NAME=shares_alert
DB_USER=your_username
DB_PASSWORD=your_password
DB_SSL_MODE=require
```

The application will automatically handle the migration!

## Email Notifications

Email notifications are sent when:
- User creates an account (welcome email)
- Price threshold alerts are triggered
- Dividend announcements (future feature)
- IPO alerts (future feature)

## Configuration Options

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `10000` |
| `DB_TYPE` | Database type (`sqlite`, `postgres`) | `sqlite` |
| `DB_FILE_PATH` | SQLite database file path | `./data/shares_alert.db` |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | Required |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret | Required |
| `JWT_SECRET` | JWT signing secret | Change in production |
| `JWT_EXPIRATION_HOURS` | JWT token expiration | `24` |
| `SMTP_HOST` | SMTP server host | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP server port | `587` |
| `SMTP_USER` | SMTP username | Required for email |
| `SMTP_PASSWORD` | SMTP password | Required for email |

## Deployment

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
```

### Environment Variables for Production
Make sure to set secure values for:
- `JWT_SECRET`
- `GOOGLE_CLIENT_SECRET`
- `SMTP_PASSWORD`
- Database credentials (if using PostgreSQL)

## Monitoring and Logging

The application includes comprehensive logging for:
- HTTP requests (via Chi middleware)
- Authentication events
- Alert processing
- Email notifications
- Database operations