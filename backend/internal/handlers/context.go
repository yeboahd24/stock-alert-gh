package handlers

import (
	"context"

	"shares-alert-backend/internal/models"
)

type contextKey string

const userContextKey contextKey = "user"

func setUserInContext(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func getUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	return user, ok
}