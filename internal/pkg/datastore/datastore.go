package datastore

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config struct holds all the configurations required the datastore package
type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`

	Username string `json:"user_name"`
	Password string `json:"password"`

	ConnPoolSize uint `json:"conn_pool_size"`
	DialTimeout  int  `json:"dial_timeout"`
}

// ConnURL returns the connection URL
func (cfg *Config) ConnURL() string {
	if cfg.Username != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	}
	return fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
}

// NewService returns a new instance of PGX pool
func NewService(cfg *Config) (*mongo.Client, error) {

	o := options.Client().ApplyURI(cfg.ConnURL())
	o.SetMaxPoolSize(uint64(cfg.ConnPoolSize))

	client, err := mongo.NewClient(o)
	if err != nil {
		return nil, fmt.Errorf("create mongo client failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DialTimeout)*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to mongo failed: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ping to mongo failed: %w", err)
	}
	return client, nil
}
