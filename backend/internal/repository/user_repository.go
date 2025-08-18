package repository

import (
	"database/sql"
	"time"

	"shares-alert-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO shares_alert_users (id, email, name, picture, google_id, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.Picture, 
		user.GoogleID, user.EmailVerified, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, name, picture, google_id, email_verified, created_at, updated_at
		FROM shares_alert_users WHERE id = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Picture,
		&user.GoogleID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, name, picture, google_id, email_verified, created_at, updated_at
		FROM shares_alert_users WHERE email = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.Picture,
		&user.GoogleID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByGoogleID(googleID string) (*models.User, error) {
	query := `
		SELECT id, email, name, picture, google_id, email_verified, created_at, updated_at
		FROM shares_alert_users WHERE google_id = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(query, googleID).Scan(
		&user.ID, &user.Email, &user.Name, &user.Picture,
		&user.GoogleID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE shares_alert_users 
		SET email = $1, name = $2, picture = $3, email_verified = $4, updated_at = $5
		WHERE id = $6
	`
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, user.Email, user.Name, user.Picture, 
		user.EmailVerified, user.UpdatedAt, user.ID)
	return err
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM shares_alert_users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// User Preferences methods
func (r *UserRepository) CreatePreferences(prefs *models.UserPreferences) error {
	query := `
		INSERT INTO shares_alert_user_preferences (id, user_id, email_notifications, push_notifications, 
			notification_frequency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, prefs.ID, prefs.UserID, prefs.EmailNotifications,
		prefs.PushNotifications, prefs.NotificationFrequency, prefs.CreatedAt, prefs.UpdatedAt)
	return err
}

func (r *UserRepository) GetPreferences(userID string) (*models.UserPreferences, error) {
	query := `
		SELECT id, user_id, email_notifications, push_notifications, 
			notification_frequency, created_at, updated_at
		FROM shares_alert_user_preferences WHERE user_id = $1
	`
	prefs := &models.UserPreferences{}
	err := r.db.QueryRow(query, userID).Scan(
		&prefs.ID, &prefs.UserID, &prefs.EmailNotifications,
		&prefs.PushNotifications, &prefs.NotificationFrequency,
		&prefs.CreatedAt, &prefs.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return prefs, nil
}

func (r *UserRepository) UpdatePreferences(prefs *models.UserPreferences) error {
	query := `
		UPDATE shares_alert_user_preferences 
		SET email_notifications = $1, push_notifications = $2, 
			notification_frequency = $3, updated_at = $4
		WHERE user_id = $5
	`
	prefs.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, prefs.EmailNotifications, prefs.PushNotifications,
		prefs.NotificationFrequency, prefs.UpdatedAt, prefs.UserID)
	return err
}