# PostgreSQL Migration Fix for Render

## Problem
User alerts were getting deleted every time the backend restarted on Render because the app was using SQLite (file-based database) on ephemeral storage.

## Solution
Migrated from SQLite to PostgreSQL using Aiven cloud database for persistent storage.

## Changes Made

### 1. Updated `render.yaml`
- Added individual database environment variables for better security
- Set `DB_TYPE=postgres` to use PostgreSQL instead of SQLite
- Configured connection to Aiven PostgreSQL with separate credentials

### 2. Updated `backend/internal/config/config.go`
- Added `parseDatabaseURL()` function to parse Render's DATABASE_URL
- Added `loadDatabaseConfig()` function to handle both DATABASE_URL and individual env vars
- Supports both PostgreSQL URL formats (`postgres://` and `postgresql://`)

## Deployment Steps

1. **Commit and push changes:**
   ```bash
   git add .
   git commit -m "Fix: Migrate from SQLite to PostgreSQL for persistent storage"
   git push origin main
   ```

2. **Deploy on Render:**
   - Render will automatically detect the new `render.yaml` configuration
   - The backend will connect to your Aiven PostgreSQL database
   - Database tables will be automatically created via migrations

3. **Verify the fix:**
   - Create some test alerts after deployment
   - Restart the backend service manually from Render dashboard
   - Check if alerts persist after restart

## Environment Variables

The backend now supports both methods:

### Method 1: Individual variables (Primary - More Secure)
```
DB_TYPE=postgres
DB_HOST=your_host
DB_PORT=5432
DB_NAME=your_database
DB_USER=your_user
DB_PASSWORD=your_password
DB_SSL_MODE=require
```

### Method 2: DATABASE_URL (Fallback for compatibility)
```
DATABASE_URL=postgres://user:password@host:port/database?sslmode=require
```

## Database Schema
The existing migration code in `database.go` will automatically create the required tables:
- `users`
- `user_preferences` 
- `alerts`

## Benefits
- ✅ Persistent storage - alerts won't be deleted on restart
- ✅ Better performance and reliability
- ✅ Automatic backups (Render PostgreSQL feature)
- ✅ Scalable for production use