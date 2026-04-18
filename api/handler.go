package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// sensitiveFields are masked in GET responses.
var sensitiveFields = map[string]bool{
	"key":            true,
	"login_password": true,
	"auth_password":  true,
}

// NewHandler returns an http.Handler that serves the /api/config endpoint.
func NewHandler(store *ConfigStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetConfig(w, store)
		case http.MethodPut:
			handlePutConfig(w, r, store)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	return mux
}

func handleGetConfig(w http.ResponseWriter, store *ConfigStore) {
	configs := store.GetAll()

	// Build response with masked sensitive fields.
	result := make(map[string]map[string]interface{}, len(configs))
	for name, cfg := range configs {
		entry := map[string]interface{}{
			"key":            cfg.Key,
			"method":         cfg.Method,
			"log_level":      cfg.LogLevel,
			"log_file":       cfg.LogFile,
			"login_password": cfg.LoginPass,
			"auth_password":  cfg.AuthPass,
			"address":        cfg.Address,
		}
		for field := range sensitiveFields {
			entry[field] = "***"
		}
		result[name] = entry
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handlePutConfig(w http.ResponseWriter, r *http.Request, store *ConfigStore) {
	var body map[string]map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	// Phase 1: Validate all services and fields before applying anything (atomicity).
	for service, fields := range body {
		if !store.HasService(service) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("service %q not found", service)})
			return
		}
		for field := range fields {
			if !IsMutable(field) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("field %q is immutable", field)})
				return
			}
		}
	}

	// Phase 2: Apply updates.
	var updated []string
	for service, fields := range body {
		stringFields := make(map[string]string, len(fields))
		for field, value := range fields {
			strVal := fmt.Sprintf("%v", value)
			stringFields[field] = strVal
			updated = append(updated, service+"."+field)

			// Side effect: hot-reload log level.
			if field == "log_level" {
				level, err := logrus.ParseLevel(strVal)
				if err == nil {
					logrus.SetLevel(level)
				}
			}
		}
		store.Update(service, stringFields)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"updated": updated})
}
