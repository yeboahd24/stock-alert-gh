package models

import (
	"time"
)

type User struct {
	ID            string    `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	Name          string    `json:"name" db:"name"`
	Picture       string    `json:"picture" db:"picture"`
	GoogleID      string    `json:"googleId" db:"google_id"`
	EmailVerified bool      `json:"emailVerified" db:"email_verified"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type UserPreferences struct {
	ID                    string `json:"id" db:"id"`
	UserID                string `json:"userId" db:"user_id"`
	EmailNotifications    bool   `json:"emailNotifications" db:"email_notifications"`
	PushNotifications     bool   `json:"pushNotifications" db:"push_notifications"`
	NotificationFrequency string `json:"notificationFrequency" db:"notification_frequency"` // immediate, daily, weekly
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`
}