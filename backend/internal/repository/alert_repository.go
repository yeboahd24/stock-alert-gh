package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"shares-alert-backend/internal/models"
)

type AlertRepository struct {
	db *sql.DB
}

func NewAlertRepository(db *sql.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) Create(alert *models.Alert) error {
	query := `
		INSERT INTO shares_alert_alerts (id, user_id, stock_symbol, stock_name, alert_type, 
			threshold_price, current_price, threshold_yield, current_yield, target_yield, 
			yield_change_threshold, last_yield, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	_, err := r.db.Exec(query, alert.ID, alert.UserID, alert.StockSymbol,
		alert.StockName, alert.AlertType, alert.ThresholdPrice, alert.CurrentPrice,
		alert.ThresholdYield, alert.CurrentYield, alert.TargetYield,
		alert.YieldChangeThreshold, alert.LastYield, alert.Status, alert.CreatedAt, alert.UpdatedAt)
	return err
}

func (r *AlertRepository) GetByID(id string) (*models.Alert, error) {
	query := `
		SELECT id, user_id, stock_symbol, stock_name, alert_type, threshold_price,
			current_price, threshold_yield, current_yield, target_yield, 
			yield_change_threshold, last_yield, status, created_at, updated_at, triggered_at
		FROM shares_alert_alerts WHERE id = $1
	`
	alert := &models.Alert{}
	err := r.db.QueryRow(query, id).Scan(
		&alert.ID, &alert.UserID, &alert.StockSymbol, &alert.StockName,
		&alert.AlertType, &alert.ThresholdPrice, &alert.CurrentPrice,
		&alert.ThresholdYield, &alert.CurrentYield, &alert.TargetYield,
		&alert.YieldChangeThreshold, &alert.LastYield, &alert.Status, 
		&alert.CreatedAt, &alert.UpdatedAt, &alert.TriggeredAt,
	)
	if err != nil {
		return nil, err
	}
	return alert, nil
}

func (r *AlertRepository) GetByUserID(userID string, filters map[string]interface{}) ([]*models.Alert, error) {
	query := `
		SELECT id, user_id, stock_symbol, stock_name, alert_type, threshold_price,
			current_price, threshold_yield, current_yield, target_yield, 
			yield_change_threshold, last_yield, status, created_at, updated_at, triggered_at
		FROM shares_alert_alerts WHERE user_id = $1
	`
	args := []interface{}{userID}
	paramCount := 1

	// Add filters
	if status, ok := filters["status"].(string); ok && status != "" {
		paramCount++
		query += fmt.Sprintf(" AND status = $%d", paramCount)
		args = append(args, status)
	}
	if stockSymbol, ok := filters["stock_symbol"].(string); ok && stockSymbol != "" {
		paramCount++
		query += fmt.Sprintf(" AND stock_symbol = $%d", paramCount)
		args = append(args, stockSymbol)
	}
	if alertType, ok := filters["alert_type"].(string); ok && alertType != "" {
		paramCount++
		query += fmt.Sprintf(" AND alert_type = $%d", paramCount)
		args = append(args, alertType)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.Alert
	for rows.Next() {
		alert := &models.Alert{}
		err := rows.Scan(
			&alert.ID, &alert.UserID, &alert.StockSymbol, &alert.StockName,
			&alert.AlertType, &alert.ThresholdPrice, &alert.CurrentPrice,
			&alert.ThresholdYield, &alert.CurrentYield, &alert.TargetYield,
			&alert.YieldChangeThreshold, &alert.LastYield, &alert.Status, 
			&alert.CreatedAt, &alert.UpdatedAt, &alert.TriggeredAt,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (r *AlertRepository) GetActiveAlerts() ([]*models.Alert, error) {
	query := `
		SELECT id, user_id, stock_symbol, stock_name, alert_type, threshold_price,
			current_price, threshold_yield, current_yield, target_yield, 
			yield_change_threshold, last_yield, status, created_at, updated_at, triggered_at
		FROM shares_alert_alerts WHERE status = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, models.AlertStatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.Alert
	for rows.Next() {
		alert := &models.Alert{}
		err := rows.Scan(
			&alert.ID, &alert.UserID, &alert.StockSymbol, &alert.StockName,
			&alert.AlertType, &alert.ThresholdPrice, &alert.CurrentPrice,
			&alert.ThresholdYield, &alert.CurrentYield, &alert.TargetYield,
			&alert.YieldChangeThreshold, &alert.LastYield, &alert.Status, 
			&alert.CreatedAt, &alert.UpdatedAt, &alert.TriggeredAt,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (r *AlertRepository) Update(alert *models.Alert) error {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}

	paramCount := 0
	
	paramCount++
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", paramCount))
	args = append(args, time.Now())

	if alert.AlertType != "" {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("alert_type = $%d", paramCount))
		args = append(args, alert.AlertType)
	}
	if alert.ThresholdPrice != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("threshold_price = $%d", paramCount))
		args = append(args, alert.ThresholdPrice)
	}
	if alert.CurrentPrice != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("current_price = $%d", paramCount))
		args = append(args, alert.CurrentPrice)
	}
	if alert.ThresholdYield != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("threshold_yield = $%d", paramCount))
		args = append(args, alert.ThresholdYield)
	}
	if alert.CurrentYield != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("current_yield = $%d", paramCount))
		args = append(args, alert.CurrentYield)
	}
	if alert.TargetYield != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("target_yield = $%d", paramCount))
		args = append(args, alert.TargetYield)
	}
	if alert.YieldChangeThreshold != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("yield_change_threshold = $%d", paramCount))
		args = append(args, alert.YieldChangeThreshold)
	}
	if alert.LastYield != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("last_yield = $%d", paramCount))
		args = append(args, alert.LastYield)
	}
	if alert.Status != "" {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("status = $%d", paramCount))
		args = append(args, alert.Status)
	}
	if alert.TriggeredAt != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("triggered_at = $%d", paramCount))
		args = append(args, alert.TriggeredAt)
	}

	args = append(args, alert.ID)

	query := fmt.Sprintf("UPDATE shares_alert_alerts SET %s WHERE id = $%d", strings.Join(setParts, ", "), len(args))
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *AlertRepository) Delete(id string) error {
	query := `DELETE FROM shares_alert_alerts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *AlertRepository) UpdateCurrentPrice(stockSymbol string, currentPrice float64) error {
	query := `
		UPDATE shares_alert_alerts 
		SET current_price = $1, updated_at = $2
		WHERE stock_symbol = $3 AND status = $4
	`
	_, err := r.db.Exec(query, currentPrice, time.Now(), stockSymbol, models.AlertStatusActive)
	return err
}

// UpdateCurrentYield updates the current dividend yield for all alerts of a specific stock
func (r *AlertRepository) UpdateCurrentYield(stockSymbol string, currentYield float64) error {
	query := `
		UPDATE shares_alert_alerts 
		SET current_yield = $1, updated_at = $2
		WHERE stock_symbol = $3 AND status = $4 AND alert_type IN ($5, $6, $7)
	`
	_, err := r.db.Exec(query, currentYield, time.Now(), stockSymbol, models.AlertStatusActive,
		models.AlertTypeHighDividendYield, models.AlertTypeDividendYieldChange, models.AlertTypeTargetDividendYield)
	return err
}

// UpdateLastYield updates the last known yield for yield change tracking
func (r *AlertRepository) UpdateLastYield(stockSymbol string, lastYield float64) error {
	query := `
		UPDATE shares_alert_alerts 
		SET last_yield = $1, updated_at = $2
		WHERE stock_symbol = $3 AND status = $4 AND alert_type = $5
	`
	_, err := r.db.Exec(query, lastYield, time.Now(), stockSymbol, models.AlertStatusActive, models.AlertTypeDividendYieldChange)
	return err
}

func (r *AlertRepository) TriggerAlert(alertID string) error {
	now := time.Now()
	query := `
		UPDATE shares_alert_alerts 
		SET status = $1, triggered_at = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(query, models.AlertStatusTriggered, now, now, alertID)
	return err
}

func (r *AlertRepository) GetActiveAlertsByType(alertType string) ([]*models.Alert, error) {
	query := `
		SELECT id, user_id, stock_symbol, stock_name, alert_type, threshold_price,
			current_price, threshold_yield, current_yield, target_yield, 
			yield_change_threshold, last_yield, status, created_at, updated_at, triggered_at
		FROM shares_alert_alerts
		WHERE status = $1 AND alert_type = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, models.AlertStatusActive, alertType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.Alert
	for rows.Next() {
		alert := &models.Alert{}
		err := rows.Scan(
			&alert.ID, &alert.UserID, &alert.StockSymbol, &alert.StockName,
			&alert.AlertType, &alert.ThresholdPrice, &alert.CurrentPrice,
			&alert.ThresholdYield, &alert.CurrentYield, &alert.TargetYield,
			&alert.YieldChangeThreshold, &alert.LastYield, &alert.Status, 
			&alert.CreatedAt, &alert.UpdatedAt, &alert.TriggeredAt,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// GetActiveDividendYieldAlerts gets all active dividend yield related alerts
func (r *AlertRepository) GetActiveDividendYieldAlerts() ([]*models.Alert, error) {
	query := `
		SELECT id, user_id, stock_symbol, stock_name, alert_type, threshold_price,
			current_price, threshold_yield, current_yield, target_yield, 
			yield_change_threshold, last_yield, status, created_at, updated_at, triggered_at
		FROM shares_alert_alerts
		WHERE status = $1 AND alert_type IN ($2, $3, $4)
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, models.AlertStatusActive, 
		models.AlertTypeHighDividendYield, models.AlertTypeDividendYieldChange, models.AlertTypeTargetDividendYield)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.Alert
	for rows.Next() {
		alert := &models.Alert{}
		err := rows.Scan(
			&alert.ID, &alert.UserID, &alert.StockSymbol, &alert.StockName,
			&alert.AlertType, &alert.ThresholdPrice, &alert.CurrentPrice,
			&alert.ThresholdYield, &alert.CurrentYield, &alert.TargetYield,
			&alert.YieldChangeThreshold, &alert.LastYield, &alert.Status, 
			&alert.CreatedAt, &alert.UpdatedAt, &alert.TriggeredAt,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}