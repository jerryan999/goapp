package configs

import (
	"time"

	"github.com/jerryan999/goapp/internal/pkg/cachestore"
	"github.com/jerryan999/goapp/internal/pkg/datastore"
	"github.com/jerryan999/goapp/internal/server/http"
)

// Configs struct handles all dependencies required for handling configurations
type Configs struct {
}

// HTTP returns the configuration required for HTTP package
func (cfg *Configs) HTTP() (*http.Config, error) {
	return &http.Config{
		Port:         "8080",
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		DialTimeout:  time.Second * 3,
	}, nil
}

// Datastore returns datastore configuration
func (cfg *Configs) Datastore() (*datastore.Config, error) {
	return &datastore.Config{
		Host:         "localhost",
		Port:         "27017",
		Username:     "",
		Password:     "",
		ConnPoolSize: 10,
		DialTimeout:  time.Second * 10,
	}, nil
}

// Cachestore returns the configuration required for cache
func (cfg *Configs) Cachestore() (*cachestore.Config, error) {
	return &cachestore.Config{
		Host: "",
		Port: "6379",

		StoreName: "0",
		Username:  "",
		Password:  "",

		PoolSize:     8,
		IdleTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		DialTimeout:  time.Second * 5,
	}, nil
}

// NewService returns an instance of Config with all the required dependencies initialized
func NewService() (*Configs, error) {
	return &Configs{}, nil
}
