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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.Picture, 
		user.GoogleID, user.EmailVerified, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, name, picture, google_id, email_verified, created_at, updated_at
		FROM shares_alert_users WHERE id = ?
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
		FROM shares_alert_users WHERE email = ?
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
		FROM shares_alert_users WHERE google_id = ?
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
		SET email = ?, name = ?, picture = ?, email_verified = ?, updated_at = ?
		WHERE id = ?
	`
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, user.Email, user.Name, user.Picture, 
		user.EmailVerified, user.UpdatedAt, user.ID)
	return err
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM shares_alert_users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// User Preferences methods
func (r *UserRepository) CreatePreferences(prefs *models.UserPreferences) error {
	query := `
		INSERT INTO user_preferences (id, user_id, email_notifications, push_notifications, 
			notification_frequency, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, prefs.ID, prefs.UserID, prefs.EmailNotifications,
		prefs.PushNotifications, prefs.NotificationFrequency, prefs.CreatedAt, prefs.UpdatedAt)
	return err
}

func (r *UserRepository) GetPreferences(userID string) (*models.UserPreferences, error) {
	query := `
		SELECT id, user_id, email_notifications, push_notifications, 
			notification_frequency, created_at, updated_at
		FROM user_preferences WHERE user_id = ?
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
		UPDATE user_preferences 
		SET email_notifications = ?, push_notifications = ?, 
			notification_frequency = ?, updated_at = ?
		WHERE user_id = ?
	`
	prefs.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, prefs.EmailNotifications, prefs.PushNotifications,
		prefs.NotificationFrequency, prefs.UpdatedAt, prefs.UserID)
	return err
}