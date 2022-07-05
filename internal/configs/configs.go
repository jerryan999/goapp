package configs

import (
	"os"
	"strconv"

	"github.com/jerryan999/goapp/internal/pkg/cachestore"
	"github.com/jerryan999/goapp/internal/pkg/datastore"
	"github.com/jerryan999/goapp/internal/server/http"
)

// AppConfigs struct handles all dependencies required for handling configurations
type AppConfigs struct {
}

// HTTP returns the configuration required for HTTP package
func (cfg *AppConfigs) HTTP() (*http.Config, error) {
	var httpConfig http.Config = http.Config{
		Port:               GetInt(os.Getenv("HTTP_PORT"), 8080),
		ReadTimeoutSecond:  GetInt(os.Getenv("HTTP_READ_TIMEOUT_SECOND"), 30),
		WriteTimeoutSecond: GetInt(os.Getenv("HTTP_WRITE_TIMEOUT_SECOND"), 30),
		DialTimeoutSecond:  GetInt(os.Getenv("HTTP_DIAL_TIMEOUT_SECOND"), 30),
	}

	return &httpConfig, nil
}

// Datastore returns datastore configuration
func (cfg *AppConfigs) Datastore() (*datastore.Config, error) {
	var dsConfig datastore.Config = datastore.Config{
		Host:         getStr(os.Getenv("DATASTORE_HOST"), "localhost"),
		Port:         GetInt(os.Getenv("DATASTORE_PORT"), 27017),
		Username:     getStr(os.Getenv("DATASTORE_USER"), ""),
		Password:     getStr(os.Getenv("DATASTORE_PASSWORD"), ""),
		ConnPoolSize: GetInt(os.Getenv("DATASTORE_CONN_POOL_SIZE"), 10),
		DialTimeout:  GetInt(os.Getenv("DATASTORE_DIAL_TIMEOUT"), 10),
	}
	return &dsConfig, nil
}

// Cachestore returns the configuration required for cache
func (cfg *AppConfigs) Cachestore() (*cachestore.Config, error) {
	var cacheConfig cachestore.Config = cachestore.Config{
		Host:         getStr(os.Getenv("CACHE_HOST"), "localhost"),
		Port:         GetInt(os.Getenv("CACHE_PORT"), 6379),
		Username:     getStr(os.Getenv("CACHE_USER"), ""),
		Password:     getStr(os.Getenv("CACHE_PASSWORD"), ""),
		PoolSize:     GetInt(os.Getenv("CACHE_POOL_SIZE"), 10),
		DialTimeout:  GetInt(os.Getenv("CACHE_DIAL_TIMEOUT"), 5),
		ReadTimeout:  GetInt(os.Getenv("CACHE_READ_TIMEOUT"), 5),
		WriteTimeout: GetInt(os.Getenv("CACHE_WRITE_TIMEOUT"), 5),
		IdleTimeout:  GetInt(os.Getenv("CACHE_IDLE_TIMEOUT"), 5),
		StoreName:    getStr(os.Getenv("CACHE_STORE_NAME"), "0"),
	}
	return &cacheConfig, nil
}

func NewService() (*AppConfigs, error) {
	cfg := AppConfigs{}
	return &cfg, nil
}

func GetInt(name string, fallback int) int {
	i, err := strconv.ParseInt(name, 10, 0)
	if nil != err {
		return fallback
	}
	return int(i)
}

func getStr(name, fallback string) string {
	if len(name) == 0 {
		return fallback
	}
	return name
}
