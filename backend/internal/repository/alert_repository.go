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
			threshold_price, current_price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query, alert.ID, alert.UserID, alert.StockSymbol,
		alert.StockName, alert.AlertType, alert.ThresholdPrice, alert.CurrentPrice,
		alert.Status, alert.CreatedAt, alert.UpdatedAt)
	return err
}

func (r *AlertRepository) GetByID(id string) (*models.Alert, error) {
	query := `
		SELECT id, user_id, stock_symbol, stock_name, alert_type, threshold_price,
			current_price, status, created_at, updated_at, triggered_at
		FROM shares_alert_alerts WHERE id = $1
	`
	alert := &models.Alert{}
	err := r.db.QueryRow(query, id).Scan(
		&alert.ID, &alert.UserID, &alert.StockSymbol, &alert.StockName,
		&alert.AlertType, &alert.ThresholdPrice, &alert.CurrentPrice,
		&alert.Status, &alert.CreatedAt, &alert.UpdatedAt, &alert.TriggeredAt,
	)
	if err != nil {
		return nil, err
	}
	return alert, nil
}

func (r *AlertRepository) GetByUserID(userID string, filters map[string]interface{}) ([]*models.Alert, error) {
	query := `
		SELECT id, user_id, stock_symbol, stock_name, alert_type, threshold_price,
			current_price, status, created_at, updated_at, triggered_at
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
			&alert.Status, &alert.CreatedAt, &alert.UpdatedAt, &alert.TriggeredAt,
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
			current_price, status, created_at, updated_at, triggered_at
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
			&alert.Status, &alert.CreatedAt, &alert.UpdatedAt, &alert.TriggeredAt,
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