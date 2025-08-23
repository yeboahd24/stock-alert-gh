-- Migration to add dividend yield fields to alerts table
-- This adds support for dividend yield alerts

-- Add new columns for dividend yield alerts
ALTER TABLE shares_alert_alerts 
ADD COLUMN threshold_yield DECIMAL(5,2),
ADD COLUMN current_yield DECIMAL(5,2),
ADD COLUMN target_yield DECIMAL(5,2),
ADD COLUMN yield_change_threshold DECIMAL(5,2),
ADD COLUMN last_yield DECIMAL(5,2);

-- Add comments for documentation
COMMENT ON COLUMN shares_alert_alerts.threshold_yield IS 'Minimum dividend yield threshold for high_dividend_yield alerts';
COMMENT ON COLUMN shares_alert_alerts.current_yield IS 'Current dividend yield of the stock';
COMMENT ON COLUMN shares_alert_alerts.target_yield IS 'Target dividend yield for target_dividend_yield alerts';
COMMENT ON COLUMN shares_alert_alerts.yield_change_threshold IS 'Minimum yield change threshold for dividend_yield_change alerts';
COMMENT ON COLUMN shares_alert_alerts.last_yield IS 'Last known dividend yield for change tracking';

-- Create index for better performance on dividend yield queries
CREATE INDEX idx_alerts_dividend_yield ON shares_alert_alerts(alert_type, current_yield) 
WHERE alert_type IN ('high_dividend_yield', 'target_dividend_yield', 'dividend_yield_change');

-- Create index for stock symbol and alert type combination
CREATE INDEX idx_alerts_stock_type ON shares_alert_alerts(stock_symbol, alert_type, status);