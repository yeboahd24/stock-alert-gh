-- Add dividend yield columns to alerts table
-- Run this SQL script directly on your database

-- For PostgreSQL:
ALTER TABLE shares_alert_alerts 
ADD COLUMN IF NOT EXISTS threshold_yield DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS current_yield DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS target_yield DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS yield_change_threshold DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS last_yield DECIMAL(5,2);

-- For SQLite (run each statement separately):
-- ALTER TABLE alerts ADD COLUMN threshold_yield DECIMAL(5,2);
-- ALTER TABLE alerts ADD COLUMN current_yield DECIMAL(5,2);
-- ALTER TABLE alerts ADD COLUMN target_yield DECIMAL(5,2);
-- ALTER TABLE alerts ADD COLUMN yield_change_threshold DECIMAL(5,2);
-- ALTER TABLE alerts ADD COLUMN last_yield DECIMAL(5,2);