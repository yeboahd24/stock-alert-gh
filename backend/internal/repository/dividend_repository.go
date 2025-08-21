package repository

import (
	"database/sql"
	"fmt"
	"time"

	"shares-alert-backend/internal/models"
)

type DividendRepository struct {
	db        *sql.DB
	tableName string
}

func NewDividendRepository(db *sql.DB, dbType string) *DividendRepository {
	tableName := "dividend_announcements"
	if dbType == "postgres" {
		tableName = "shares_alert_dividend_announcements"
	}
	return &DividendRepository{
		db:        db,
		tableName: tableName,
	}
}

func (r *DividendRepository) Create(dividend *models.DividendAnnouncement) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, stock_symbol, stock_name, dividend_type, amount, currency, ex_date, payment_date, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.tableName)
	_, err := r.db.Exec(query, dividend.ID, dividend.StockSymbol, dividend.StockName, dividend.DividendType, dividend.Amount, dividend.Currency, dividend.ExDate, dividend.PaymentDate, dividend.Status, dividend.CreatedAt, dividend.UpdatedAt)
	return err
}

func (r *DividendRepository) GetAll() ([]*models.DividendAnnouncement, error) {
	query := fmt.Sprintf(`
		SELECT id, stock_symbol, stock_name, dividend_type, amount, currency, ex_date, payment_date, status, created_at, updated_at
		FROM %s
		ORDER BY ex_date DESC
	`, r.tableName)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dividends []*models.DividendAnnouncement
	for rows.Next() {
		dividend := &models.DividendAnnouncement{}
		err := rows.Scan(&dividend.ID, &dividend.StockSymbol, &dividend.StockName, &dividend.DividendType, &dividend.Amount, &dividend.Currency, &dividend.ExDate, &dividend.PaymentDate, &dividend.Status, &dividend.CreatedAt, &dividend.UpdatedAt)
		if err != nil {
			return nil, err
		}
		dividends = append(dividends, dividend)
	}
	return dividends, nil
}

func (r *DividendRepository) GetUpcoming() ([]*models.DividendAnnouncement, error) {
	query := fmt.Sprintf(`
		SELECT id, stock_symbol, stock_name, dividend_type, amount, currency, ex_date, payment_date, status, created_at, updated_at
		FROM %s
		WHERE status = ? AND ex_date > ?
		ORDER BY ex_date ASC
	`, r.tableName)
	rows, err := r.db.Query(query, models.DividendStatusAnnounced, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dividends []*models.DividendAnnouncement
	for rows.Next() {
		dividend := &models.DividendAnnouncement{}
		err := rows.Scan(&dividend.ID, &dividend.StockSymbol, &dividend.StockName, &dividend.DividendType, &dividend.Amount, &dividend.Currency, &dividend.ExDate, &dividend.PaymentDate, &dividend.Status, &dividend.CreatedAt, &dividend.UpdatedAt)
		if err != nil {
			return nil, err
		}
		dividends = append(dividends, dividend)
	}
	return dividends, nil
}

func (r *DividendRepository) UpdateStatus(id, status string) error {
	query := fmt.Sprintf(`UPDATE %s SET status = ?, updated_at = ? WHERE id = ?`, r.tableName)
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *DividendRepository) GetBySymbol(symbol string) ([]*models.DividendAnnouncement, error) {
	query := fmt.Sprintf(`
		SELECT id, stock_symbol, stock_name, dividend_type, amount, currency, ex_date, payment_date, status, created_at, updated_at
		FROM %s
		WHERE stock_symbol = ?
		ORDER BY ex_date DESC
	`, r.tableName)
	rows, err := r.db.Query(query, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dividends []*models.DividendAnnouncement
	for rows.Next() {
		dividend := &models.DividendAnnouncement{}
		err := rows.Scan(&dividend.ID, &dividend.StockSymbol, &dividend.StockName, &dividend.DividendType, &dividend.Amount, &dividend.Currency, &dividend.ExDate, &dividend.PaymentDate, &dividend.Status, &dividend.CreatedAt, &dividend.UpdatedAt)
		if err != nil {
			return nil, err
		}
		dividends = append(dividends, dividend)
	}
	return dividends, nil
}