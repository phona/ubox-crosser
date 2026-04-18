package api

import "net/http"

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
// TODO: implement — currently returns 501 so contract tests fail.
func NewRoutesHandler(provider RouteProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}
