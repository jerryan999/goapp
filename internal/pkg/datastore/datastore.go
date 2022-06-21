package datastore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config struct holds all the configurations required the datastore package
type Config struct {
	Host string `json:"host,omitempty"`
	Port string `json:"port,omitempty"`

	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`

	ConnPoolSize uint          `json:"connPoolSize,omitempty"`
	DialTimeout  time.Duration `json:"dialTimeout,omitempty"`
}

// ConnURL returns the connection URL
func (cfg *Config) ConnURL() string {
	if cfg.Username != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	}
	return fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)
}

// NewService returns a new instance of PGX pool
func NewService(cfg *Config) (*mongo.Client, error) {

	o := options.Client().ApplyURI(cfg.ConnURL())
	o.SetMaxPoolSize(uint64(cfg.ConnPoolSize))

	client, err := mongo.NewClient(o)
	if err != nil {
		return nil, errors.New("create mongo client failed: %w")
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.New("connect to mongo failed: %w")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.New("ping to mongo failed: %w")
	}
	return client, nil
}
