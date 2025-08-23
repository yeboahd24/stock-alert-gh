# 🚨 PRODUCTION HOTFIX GUIDE

## Issue
**Error**: `pq: column "threshold_yield" does not exist`

**Root Cause**: The dividend yield migration was not applied to the production database, but the application code is trying to query these columns.

## 🚀 Quick Fix Options

### Option 1: Automated Fix (Recommended)
```bash
./fix_production_deployment.sh
```

This script will:
- Build and run the migration tool
- Add missing columns to the database
- Test the build
- Commit and deploy the fix

### Option 2: Manual Database Fix
If you have direct database access, run this SQL:

```sql
-- Add missing columns
ALTER TABLE shares_alert_alerts 
ADD COLUMN IF NOT EXISTS threshold_yield DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS current_yield DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS target_yield DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS yield_change_threshold DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS last_yield DECIMAL(5,2);

-- Add performance indexes
CREATE INDEX IF NOT EXISTS idx_alerts_dividend_yield ON shares_alert_alerts(alert_type, current_yield) 
WHERE alert_type IN ('high_dividend_yield', 'target_dividend_yield', 'dividend_yield_change');

CREATE INDEX IF NOT EXISTS idx_alerts_stock_type ON shares_alert_alerts(stock_symbol, alert_type, status);
```

### Option 3: Using Migration Tool Directly
```bash
cd backend
go run ./cmd/migrate
```

## 🔍 Verification

After applying the fix, verify with:

```sql
-- Check if columns exist
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'shares_alert_alerts' 
AND column_name IN ('threshold_yield', 'current_yield', 'target_yield', 'yield_change_threshold', 'last_yield');
```

Expected output: 5 rows showing the new columns.

## 🧪 Testing

1. **Health Check**: `GET /api/v1/health`
2. **Alerts Endpoint**: `GET /api/v1/alerts`
3. **Check Logs**: Look for any remaining database errors

## 📊 Impact

- **Downtime**: Minimal (migration runs quickly)
- **Data Loss**: None (only adding columns)
- **Compatibility**: Backward compatible (uses IF NOT EXISTS)

## 🔄 Prevention

To prevent this in the future:
1. Always run migrations before deploying code changes
2. Include migration steps in deployment pipeline
3. Test migrations in staging environment first

## 📞 Emergency Contacts

If the automated fix fails:
1. Check Render deployment logs
2. Verify database connection
3. Run manual SQL script
4. Contact database administrator if needed

---

**Status**: Ready to deploy
**Priority**: Critical (production down)
**ETA**: 5-10 minutes