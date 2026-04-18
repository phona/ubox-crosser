package api

import (
	"sync"

	"github.com/phona/ubox-crosser/models/config"
)

// mutableFields defines which fields can be updated at runtime.
var mutableFields = map[string]bool{
	"log_level": true,
	"log_file":  true,
}

// ConfigStore provides thread-safe access to service configurations.
type ConfigStore struct {
	mu      sync.RWMutex
	configs map[string]config.ServerConfig
}

// NewConfigStore creates a ConfigStore from a map of service configs.
func NewConfigStore(configs map[string]config.ServerConfig) *ConfigStore {
	// Copy the map to avoid external mutation.
	m := make(map[string]config.ServerConfig, len(configs))
	for k, v := range configs {
		m[k] = v
	}
	return &ConfigStore{configs: m}
}

// GetAll returns a snapshot copy of all service configurations.
func (s *ConfigStore) GetAll() map[string]config.ServerConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]config.ServerConfig, len(s.configs))
	for k, v := range s.configs {
		result[k] = v
	}
	return result
}

// IsMutable returns true if the given field name can be updated at runtime.
func IsMutable(field string) bool {
	return mutableFields[field]
}

// Update atomically applies mutable field updates for a single service.
// It modifies the config in place under the write lock.
func (s *ConfigStore) Update(service string, fields map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg := s.configs[service]
	for field, value := range fields {
		switch field {
		case "log_level":
			cfg.LogLevel = value
		case "log_file":
			cfg.LogFile = value
		}
	}
	s.configs[service] = cfg
}

// HasService returns true if the service exists in the store.
func (s *ConfigStore) HasService(service string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.configs[service]
	return ok
}
