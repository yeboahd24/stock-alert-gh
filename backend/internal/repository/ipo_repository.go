package repository

import (
	"database/sql"
	"fmt"
	"time"

	"shares-alert-backend/internal/models"
)

type IPORepository struct {
	db        *sql.DB
	tableName string
}

func NewIPORepository(db *sql.DB, dbType string) *IPORepository {
	tableName := "ipo_announcements"
	if dbType == "postgres" {
		tableName = "shares_alert_ipo_announcements"
	}
	return &IPORepository{
		db:        db,
		tableName: tableName,
	}
}

func (r *IPORepository) Create(ipo *models.IPOAnnouncement) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, company_name, symbol, sector, offer_price, listing_date, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.tableName)
	_, err := r.db.Exec(query, ipo.ID, ipo.CompanyName, ipo.Symbol, ipo.Sector, ipo.OfferPrice, ipo.ListingDate, ipo.Status, ipo.CreatedAt, ipo.UpdatedAt)
	return err
}

func (r *IPORepository) GetAll() ([]*models.IPOAnnouncement, error) {
	query := fmt.Sprintf(`
		SELECT id, company_name, symbol, sector, offer_price, listing_date, status, created_at, updated_at
		FROM %s
		ORDER BY listing_date DESC
	`, r.tableName)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ipos []*models.IPOAnnouncement
	for rows.Next() {
		ipo := &models.IPOAnnouncement{}
		err := rows.Scan(&ipo.ID, &ipo.CompanyName, &ipo.Symbol, &ipo.Sector, &ipo.OfferPrice, &ipo.ListingDate, &ipo.Status, &ipo.CreatedAt, &ipo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		ipos = append(ipos, ipo)
	}
	return ipos, nil
}

func (r *IPORepository) GetUpcoming() ([]*models.IPOAnnouncement, error) {
	query := fmt.Sprintf(`
		SELECT id, company_name, symbol, sector, offer_price, listing_date, status, created_at, updated_at
		FROM %s
		WHERE status = ? AND listing_date > ?
		ORDER BY listing_date ASC
	`, r.tableName)
	rows, err := r.db.Query(query, models.IPOStatusAnnounced, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ipos []*models.IPOAnnouncement
	for rows.Next() {
		ipo := &models.IPOAnnouncement{}
		err := rows.Scan(&ipo.ID, &ipo.CompanyName, &ipo.Symbol, &ipo.Sector, &ipo.OfferPrice, &ipo.ListingDate, &ipo.Status, &ipo.CreatedAt, &ipo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		ipos = append(ipos, ipo)
	}
	return ipos, nil
}

func (r *IPORepository) UpdateStatus(id, status string) error {
	query := fmt.Sprintf(`UPDATE %s SET status = ?, updated_at = ? WHERE id = ?`, r.tableName)
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *IPORepository) GetBySymbol(symbol string) (*models.IPOAnnouncement, error) {
	query := fmt.Sprintf(`
		SELECT id, company_name, symbol, sector, offer_price, listing_date, status, created_at, updated_at
		FROM %s
		WHERE symbol = ?
	`, r.tableName)
	ipo := &models.IPOAnnouncement{}
	err := r.db.QueryRow(query, symbol).Scan(&ipo.ID, &ipo.CompanyName, &ipo.Symbol, &ipo.Sector, &ipo.OfferPrice, &ipo.ListingDate, &ipo.Status, &ipo.CreatedAt, &ipo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return ipo, nil
}