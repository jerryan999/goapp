package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jerryan999/goapp/internal/pkg/cachestore"
	"github.com/jerryan999/goapp/internal/pkg/datastore"
	"github.com/jerryan999/goapp/internal/server/http"
)

// AppConfigs struct handles all dependencies required for handling configurations
type AppConfigs struct {
	data map[string]json.RawMessage
}

// HTTP returns the configuration required for HTTP package
func (cfg *AppConfigs) HTTP() (*http.Config, error) {
	var httpConfig http.Config
	err := json.Unmarshal(cfg.data["http"], &httpConfig)
	if err != nil {
		return nil, fmt.Errorf("http config: %w", err)
	}
	return &httpConfig, nil
}

// Datastore returns datastore configuration
func (cfg *AppConfigs) Datastore() (*datastore.Config, error) {
	var dsConfig datastore.Config
	err := json.Unmarshal(cfg.data["mongo"], &dsConfig)
	if err != nil {
		return nil, fmt.Errorf("datastore config: %w", err)
	}
	return &dsConfig, nil
}

// Cachestore returns the configuration required for cache
func (cfg *AppConfigs) Cachestore() (*cachestore.Config, error) {
	var cacheConfig cachestore.Config
	err := json.Unmarshal(cfg.data["redis"], &cacheConfig)
	if err != nil {
		return nil, fmt.Errorf("cachestore config: %w", err)
	}
	return &cacheConfig, nil
}

func (cfg *AppConfigs) loadJsonConfig(yamlFilePath string) error {
	content, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &cfg.data)
	return err
}

func NewService(jsonFilePath string) (*AppConfigs, error) {
	var cfg AppConfigs
	err := cfg.loadJsonConfig(jsonFilePath)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
