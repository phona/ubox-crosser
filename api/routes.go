package api

import (
	"encoding/json"
	"net/http"
)

// RouteInfo represents a single proxy route as defined in contract.spec.yaml.
type RouteInfo struct {
	Name          string `json:"name"`
	ListenAddress string `json:"listen_address"`
	Method        string `json:"method"`
	Active        bool   `json:"active"`
}

// RouteProvider supplies the list of configured routes.
type RouteProvider interface {
	ListRoutes() []RouteInfo
}

// NewRoutesHandler returns an http.Handler for GET /api/routes.
func NewRoutesHandler(provider RouteProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var routes []RouteInfo
		if provider != nil {
			routes = provider.ListRoutes()
		}
		if routes == nil {
			routes = []RouteInfo{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes)
	})
}
