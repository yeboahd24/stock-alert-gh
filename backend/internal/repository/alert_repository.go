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
		INSERT INTO alerts (id, user_id, stock_symbol, stock_name, alert_type, 
			threshold_price, current_price, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		FROM alerts WHERE id = ?
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
		FROM alerts WHERE user_id = ?
	`
	args := []interface{}{userID}

	// Add filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if stockSymbol, ok := filters["stock_symbol"].(string); ok && stockSymbol != "" {
		query += " AND stock_symbol = ?"
		args = append(args, stockSymbol)
	}
	if alertType, ok := filters["alert_type"].(string); ok && alertType != "" {
		query += " AND alert_type = ?"
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
		FROM alerts WHERE status = ?
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

	setParts = append(setParts, "updated_at = ?")
	args = append(args, time.Now())

	if alert.AlertType != "" {
		setParts = append(setParts, "alert_type = ?")
		args = append(args, alert.AlertType)
	}
	if alert.ThresholdPrice != nil {
		setParts = append(setParts, "threshold_price = ?")
		args = append(args, alert.ThresholdPrice)
	}
	if alert.CurrentPrice != nil {
		setParts = append(setParts, "current_price = ?")
		args = append(args, alert.CurrentPrice)
	}
	if alert.Status != "" {
		setParts = append(setParts, "status = ?")
		args = append(args, alert.Status)
	}
	if alert.TriggeredAt != nil {
		setParts = append(setParts, "triggered_at = ?")
		args = append(args, alert.TriggeredAt)
	}

	args = append(args, alert.ID)

	query := fmt.Sprintf("UPDATE alerts SET %s WHERE id = ?", strings.Join(setParts, ", "))
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *AlertRepository) Delete(id string) error {
	query := `DELETE FROM alerts WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *AlertRepository) UpdateCurrentPrice(stockSymbol string, currentPrice float64) error {
	query := `
		UPDATE alerts 
		SET current_price = ?, updated_at = ?
		WHERE stock_symbol = ? AND status = ?
	`
	_, err := r.db.Exec(query, currentPrice, time.Now(), stockSymbol, models.AlertStatusActive)
	return err
}

func (r *AlertRepository) TriggerAlert(alertID string) error {
	now := time.Now()
	query := `
		UPDATE alerts 
		SET status = ?, triggered_at = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, models.AlertStatusTriggered, now, now, alertID)
	return err
}