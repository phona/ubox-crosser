package api

import (
	"fmt"
	"sync"

	"github.com/phona/ubox-crosser/models/config"
)

// mutableFields lists config fields that can be updated at runtime.
var mutableFields = map[string]bool{
	"log_level": true,
	"log_file":  true,
}

// ConfigStore provides thread-safe access to service configurations.
type ConfigStore struct {
	mu      sync.RWMutex
	configs map[string]config.ServerConfig
}

// NewConfigStore creates a ConfigStore with the given initial configurations.
func NewConfigStore(configs map[string]config.ServerConfig) *ConfigStore {
	return &ConfigStore{configs: configs}
}

// GetAll returns a copy of all service configurations.
func (s *ConfigStore) GetAll() map[string]config.ServerConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]config.ServerConfig, len(s.configs))
	for k, v := range s.configs {
		result[k] = v
	}
	return result
}

// Update validates and applies mutable field updates atomically.
// Returns the list of updated field paths (e.g. "service.log_level").
func (s *ConfigStore) Update(updates map[string]map[string]string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate all updates before applying any (atomicity).
	for svcName, fields := range updates {
		if _, ok := s.configs[svcName]; !ok {
			return nil, &ServiceNotFoundError{Service: svcName}
		}
		for field := range fields {
			if !mutableFields[field] {
				return nil, &ImmutableFieldError{Field: field}
			}
		}
	}

	// Apply updates.
	var updated []string
	for svcName, fields := range updates {
		cfg := s.configs[svcName]
		for field, value := range fields {
			switch field {
			case "log_level":
				cfg.LogLevel = value
			case "log_file":
				cfg.LogFile = value
			}
			updated = append(updated, fmt.Sprintf("%s.%s", svcName, field))
		}
		s.configs[svcName] = cfg
	}

	return updated, nil
}

// ServiceNotFoundError indicates the requested service does not exist.
type ServiceNotFoundError struct {
	Service string
}

func (e *ServiceNotFoundError) Error() string {
	return fmt.Sprintf("service %q not found", e.Service)
}

// ImmutableFieldError indicates an attempt to update an immutable field.
type ImmutableFieldError struct {
	Field string
}

func (e *ImmutableFieldError) Error() string {
	return fmt.Sprintf("field %q is immutable", e.Field)
}
