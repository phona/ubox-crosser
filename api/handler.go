package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"

	"github.com/sirupsen/logrus"
)

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

	// Build response with sensitive fields masked.
	resp := make(map[string]map[string]interface{}, len(configs))
	for name, cfg := range configs {
		resp[name] = map[string]interface{}{
			"key":            "***",
			"method":         cfg.Method,
			"log_file":       cfg.LogFile,
			"log_level":      cfg.LogLevel,
			"config_file":    cfg.ConfigFile,
			"login_password": "***",
			"auth_password":  "***",
			"address":        cfg.Address,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handlePutConfig(w http.ResponseWriter, r *http.Request, store *ConfigStore) {
	var body map[string]map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	updated, err := store.Update(body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		var svcErr *ServiceNotFoundError
		if errors.As(err, &svcErr) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Apply hot-reload side effects.
	for _, path := range updated {
		// path format: "service.field"
		if len(path) > 10 && path[len(path)-10:] == ".log_level" {
			configs := store.GetAll()
			for _, cfg := range configs {
				level, parseErr := logrus.ParseLevel(cfg.LogLevel)
				if parseErr == nil {
					logrus.SetLevel(level)
				}
			}
			break
		}
	}

	sort.Strings(updated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"updated": updated})
}
