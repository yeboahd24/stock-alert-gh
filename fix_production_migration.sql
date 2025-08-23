-- Emergency production fix for missing dividend yield columns
-- This script adds the missing columns that are causing the production error

-- Add dividend yield columns to alerts table
ALTER TABLE shares_alert_alerts 
ADD COLUMN IF NOT EXISTS threshold_yield DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS current_yield DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS target_yield DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS yield_change_threshold DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS last_yield DECIMAL(5,2);

-- Add comments for documentation
COMMENT ON COLUMN shares_alert_alerts.threshold_yield IS 'Minimum dividend yield threshold for high_dividend_yield alerts';
COMMENT ON COLUMN shares_alert_alerts.current_yield IS 'Current dividend yield of the stock';
COMMENT ON COLUMN shares_alert_alerts.target_yield IS 'Target dividend yield for target_dividend_yield alerts';
COMMENT ON COLUMN shares_alert_alerts.yield_change_threshold IS 'Minimum yield change threshold for dividend_yield_change alerts';
COMMENT ON COLUMN shares_alert_alerts.last_yield IS 'Last known dividend yield for change tracking';

-- Create index for better performance on dividend yield queries
CREATE INDEX IF NOT EXISTS idx_alerts_dividend_yield ON shares_alert_alerts(alert_type, current_yield) 
WHERE alert_type IN ('high_dividend_yield', 'target_dividend_yield', 'dividend_yield_change');

-- Create index for stock symbol and alert type combination
CREATE INDEX IF NOT EXISTS idx_alerts_stock_type ON shares_alert_alerts(stock_symbol, alert_type, status);

-- Verify the columns were added successfully
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'shares_alert_alerts' 
AND column_name IN ('threshold_yield', 'current_yield', 'target_yield', 'yield_change_threshold', 'last_yield')
ORDER BY column_name;