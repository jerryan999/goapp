package users

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jerryan999/goapp/internal/pkg/cachestore"
	"github.com/jerryan999/goapp/internal/pkg/logger"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserValidation = errors.New("validation error")
)

// User holds all data required to represent a user
type User struct {
	FirstName string     `json:"firstName,omitempty"`
	LastName  string     `json:"lastName,omitempty"`
	Mobile    string     `json:"mobile,omitempty"`
	Email     string     `json:"email,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (u *User) setDefaults() {
	now := time.Now()
	if u.CreatedAt == nil {
		u.CreatedAt = &now
	}

	if u.UpdatedAt == nil {
		u.UpdatedAt = &now
	}
}

// Sanitize is used to sanitize/cleanup the fields of User
func (u *User) Sanitize() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
	u.Mobile = strings.TrimSpace(u.Mobile)
}

// Validate is used to validate the fields of User
func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("validate:%w", ErrUserValidation)
	}

	err := validateEmail(u.Email)
	if err != nil {
		return err
	}

	return nil
}

func validateEmail(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("validateEmail: %w, email: %s", ErrUserValidation, email)
	}

	return nil
}

type Users struct {
	logHandler logger.Logger
	cachestore userCachestore
	store      store
}

// CreateUser creates a new user
func (us *Users) CreateUser(ctx context.Context, u *User) (*User, error) {
	u.setDefaults()
	u.Sanitize()

	err := u.Validate()
	if err != nil {
		if errors.Is(err, ErrUserValidation) {
			us.logHandler.Warn(err.Error())
		} else {
			us.logHandler.Error(err.Error())
		}
		return nil, err
	}

	err = us.store.Create(ctx, u)
	if err != nil {
		us.logHandler.Error(err.Error())
		return nil, err
	}

	return u, nil
}

// ReadByEmail returns a user which matches the given email
func (us *Users) ReadByEmail(ctx context.Context, email string) (*User, error) {
	email = strings.TrimSpace(email)
	err := validateEmail(email)
	if err != nil {
		us.logHandler.Info("ReadByEmail: %s", err.Error())
		return nil, err
	}

	u, err := us.cachestore.ReadUserByEmail(ctx, email)
	if err != nil &&
		!errors.Is(err, cachestore.ErrCacheMiss) &&
		!errors.Is(err, cachestore.ErrCacheNotInitialized) {
		// caches are usually read-through, i.e. in case of error, just log and continue to fetch from
		// primary datastore
		us.logHandler.Error(err.Error())
	} else if err == nil {
		return u, nil
	}

	u, err = us.store.ReadByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("readByEmail: %w", ErrUserNotFound)
	}

	err = us.cachestore.SetUser(ctx, u.Email, u)
	if err != nil {
		// in case of error while storing in cache, it is only logged
		// This behaviour as well as read-through cache behaviour depends on your business logic.
		us.logHandler.Error(err.Error())
	}

	return u, nil
}

// NewService initializes the Users struct with all its dependencies and returns a new instance
// all dependencies of Users should be sent as arguments of NewService
func NewService(l logger.Logger, m *mongo.Client, redispool *redis.Pool) (*Users, error) {
	ustore, err := newStore(m)
	if err != nil {
		return nil, err
	}

	cstore, err := newCacheStore(redispool)
	if err != nil {
		return nil, err
	}

	return &Users{
		logHandler: l,
		cachestore: cstore,
		store:      ustore,
	}, nil
}
