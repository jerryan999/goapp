package api

import (
	"context"

	"github.com/jerryan999/goapp/internal/users"
)

// CreateUser is the API to create/signup a new user
func (a *API) CreateUser(ctx context.Context, u *users.User) (*users.User, error) {
	return a.users.CreateUser(ctx, u)
}

// ReadUserByEmail is the API to read an existing user by their email
func (a *API) ReadUserByEmail(ctx context.Context, email string) (*users.User, error) {
	return a.users.ReadByEmail(ctx, email)
}
