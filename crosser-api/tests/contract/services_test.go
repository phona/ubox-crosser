package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func stubServicesMux() *http.ServeMux {
	resetStubServices()
	now := time.Now().Format(time.RFC3339)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/services", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			svcs := make([]Service, 0, len(stubServices))
			for _, s := range stubServices {
				svcs = append(svcs, s)
			}
			writeSuccess(w, http.StatusOK, ServiceListData{Services: svcs, Total: len(svcs)})
		case http.MethodPost:
			var req struct {
				Name          string `json:"name"`
				Key           string `json:"key"`
				Method        string `json:"method"`
				Address       string `json:"address"`
				ExposeAddress string `json:"expose_address"`
				LoginPassword string `json:"login_password"`
				AuthPassword  string `json:"auth_password"`
				LogFile       string `json:"log_file"`
				LogLevel      string `json:"log_level"`
			}
			if err := readJSON(r, &req); err != nil {
				writeError(w, http.StatusBadRequest, 2003, "invalid body")
				return
			}
			if _, exists := stubServices[req.Name]; exists {
				writeError(w, http.StatusConflict, 2002, "service already exists")
				return
			}
			svc := Service{
				Name:          req.Name,
				Key:           req.Key,
				Method:        req.Method,
				Address:       req.Address,
				ExposeAddress: req.ExposeAddress,
				LoginPassword: req.LoginPassword,
				AuthPassword:  req.AuthPassword,
				LogFile:       req.LogFile,
				LogLevel:      req.LogLevel,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			stubServices[req.Name] = svc
			writeSuccess(w, http.StatusCreated, svc)
		default:
			writeError(w, http.StatusMethodNotAllowed, 9999, "method not allowed")
		}
	})
	mux.HandleFunc("/api/v1/services/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/services/")
		parts := strings.SplitN(path, "/", 2)
		name := parts[0]

		if len(parts) == 2 && parts[1] == "config" {
			svc, ok := stubServices[name]
			if !ok {
				writeError(w, http.StatusNotFound, 2001, "service not found")
				return
			}
			config := map[string]interface{}{
				"common": map[string]string{
					"key":            svc.Key,
					"method":         svc.Method,
					"address":        svc.Address,
					"login_password": svc.LoginPassword,
					"auth_password":  svc.AuthPassword,
					"log_file":       svc.LogFile,
					"log_level":      svc.LogLevel,
				},
				name: map[string]string{"key": svc.Key},
			}
			writeSuccess(w, http.StatusOK, config)
			return
		}

		switch r.Method {
		case http.MethodGet:
			svc, ok := stubServices[name]
			if !ok {
				writeError(w, http.StatusNotFound, 2001, "service not found")
				return
			}
			writeSuccess(w, http.StatusOK, ServiceDetailData{
				Service:        svc,
				ProxyInstances: []ProxyInstance{},
			})
		case http.MethodPut:
			svc, ok := stubServices[name]
			if !ok {
				writeError(w, http.StatusNotFound, 2001, "service not found")
				return
			}
			var req map[string]string
			if err := readJSON(r, &req); err != nil {
				writeError(w, http.StatusBadRequest, 2003, "invalid body")
				return
			}
			if v, ok := req["key"]; ok {
				svc.Key = v
			}
			if v, ok := req["method"]; ok {
				svc.Method = v
			}
			if v, ok := req["address"]; ok {
				svc.Address = v
			}
			svc.UpdatedAt = time.Now().Add(time.Second).Format(time.RFC3339)
			stubServices[name] = svc
			writeSuccess(w, http.StatusOK, svc)
		case http.MethodDelete:
			if _, ok := stubServices[name]; !ok {
				writeError(w, http.StatusNotFound, 2001, "service not found")
				return
			}
			delete(stubServices, name)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			writeJSON(w, map[string]interface{}{"code": 0, "message": "success"})
		default:
			writeError(w, http.StatusMethodNotAllowed, 9999, "method not allowed")
		}
	})
	return mux
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	_ = json.NewEncoder(w).Encode(v)
}

// REQ-997-S4: POST /api/v1/services creates a service and returns 201
func TestCreateService(t *testing.T) {
	srv := httptest.NewServer(stubServicesMux())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodPost, "/api/v1/services",
		jsonBody(t, map[string]string{
			"name": "test-svc", "key": "k1", "method": "aes-256-cfb", "address": ":8388",
		}), nil)

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 0 {
		t.Fatalf("expected code 0, got %d", ar.Code)
	}

	svc := parseData[Service](t, ar)
	if svc.Name != "test-svc" {
		t.Errorf("expected name test-svc, got %q", svc.Name)
	}
	if svc.CreatedAt == "" {
		t.Error("expected non-empty created_at")
	}
	if _, err := time.Parse(time.RFC3339, svc.CreatedAt); err != nil {
		t.Errorf("created_at not RFC3339: %v", err)
	}
}

// REQ-997-S5: GET /api/v1/services returns service list
func TestListServices(t *testing.T) {
	srv := httptest.NewServer(stubServicesMux())
	defer srv.Close()

	doRequest(t, srv, http.MethodPost, "/api/v1/services",
		jsonBody(t, map[string]string{
			"name": "svc-a", "key": "k1", "method": "aes-256-cfb", "address": ":8388",
		}), nil)

	resp := doRequest(t, srv, http.MethodGet, "/api/v1/services", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 0 {
		t.Fatalf("expected code 0, got %d", ar.Code)
	}

	list := parseData[ServiceListData](t, ar)
	if len(list.Services) == 0 {
		t.Error("expected non-empty services array")
	}
	if list.Total <= 0 {
		t.Errorf("expected positive total, got %d", list.Total)
	}
	for _, s := range list.Services {
		if s.Name == "" || s.Key == "" || s.Method == "" || s.Address == "" {
			t.Errorf("service missing required fields: %+v", s)
		}
	}
}

// REQ-997-S6: GET /api/v1/services/{name} returns service details
func TestGetServiceDetail(t *testing.T) {
	srv := httptest.NewServer(stubServicesMux())
	defer srv.Close()

	doRequest(t, srv, http.MethodPost, "/api/v1/services",
		jsonBody(t, map[string]string{
			"name": "test-svc", "key": "k1", "method": "aes-256-cfb", "address": ":8388",
		}), nil)

	resp := doRequest(t, srv, http.MethodGet, "/api/v1/services/test-svc", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 0 {
		t.Fatalf("expected code 0, got %d", ar.Code)
	}

	detail := parseData[ServiceDetailData](t, ar)
	if detail.Service.Name != "test-svc" {
		t.Errorf("expected service name test-svc, got %q", detail.Service.Name)
	}
	if detail.ProxyInstances == nil {
		t.Error("expected proxy_instances array (even if empty)")
	}
}

// REQ-997-S7: PUT /api/v1/services/{name} updates service
func TestUpdateService(t *testing.T) {
	srv := httptest.NewServer(stubServicesMux())
	defer srv.Close()

	createResp := doRequest(t, srv, http.MethodPost, "/api/v1/services",
		jsonBody(t, map[string]string{
			"name": "test-svc", "key": "k1", "method": "aes-256-cfb", "address": ":8388",
		}), nil)
	createAr := parseResponse(t, createResp)
	origSvc := parseData[Service](t, createAr)

	resp := doRequest(t, srv, http.MethodPut, "/api/v1/services/test-svc",
		jsonBody(t, map[string]string{"key": "new-key"}), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 0 {
		t.Fatalf("expected code 0, got %d", ar.Code)
	}

	updated := parseData[Service](t, ar)
	if updated.Key != "new-key" {
		t.Errorf("expected key new-key, got %q", updated.Key)
	}
	if updated.UpdatedAt == origSvc.CreatedAt {
		t.Error("expected updated_at to differ from created_at")
	}
}

// REQ-997-S8: DELETE /api/v1/services/{name} deletes a service
func TestDeleteService(t *testing.T) {
	srv := httptest.NewServer(stubServicesMux())
	defer srv.Close()

	doRequest(t, srv, http.MethodPost, "/api/v1/services",
		jsonBody(t, map[string]string{
			"name": "test-svc", "key": "k1", "method": "aes-256-cfb", "address": ":8388",
		}), nil)

	resp := doRequest(t, srv, http.MethodDelete, "/api/v1/services/test-svc", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	getResp := doRequest(t, srv, http.MethodGet, "/api/v1/services/test-svc", nil, nil)
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}

// REQ-997-S9: GET non-existent service returns 404
func TestGetServiceNotFound(t *testing.T) {
	srv := httptest.NewServer(stubServicesMux())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/v1/services/no-such-svc", nil, nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 2001 {
		t.Errorf("expected code 2001, got %d", ar.Code)
	}
}

// REQ-997-S10: POST duplicate service returns 409
func TestCreateDuplicateService(t *testing.T) {
	srv := httptest.NewServer(stubServicesMux())
	defer srv.Close()

	body := jsonBody(t, map[string]string{
		"name": "test-svc", "key": "k1", "method": "aes-256-cfb", "address": ":8388",
	})
	doRequest(t, srv, http.MethodPost, "/api/v1/services", body, nil)

	body2 := jsonBody(t, map[string]string{
		"name": "test-svc", "key": "k2", "method": "aes-256-cfb", "address": ":8389",
	})
	resp := doRequest(t, srv, http.MethodPost, "/api/v1/services", body2, nil)
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 2002 {
		t.Errorf("expected code 2002, got %d", ar.Code)
	}
}

// REQ-997-S11: GET /api/v1/services/{name}/config returns server.json compatible format
func TestGetServiceConfig(t *testing.T) {
	srv := httptest.NewServer(stubServicesMux())
	defer srv.Close()

	doRequest(t, srv, http.MethodPost, "/api/v1/services",
		jsonBody(t, map[string]string{
			"name": "test-svc", "key": "k1", "method": "aes-256-cfb", "address": ":8388",
			"login_password": "lp", "auth_password": "ap", "log_file": "", "log_level": "info",
		}), nil)

	resp := doRequest(t, srv, http.MethodGet, "/api/v1/services/test-svc/config", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 0 {
		t.Fatalf("expected code 0, got %d", ar.Code)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(ar.Data, &config); err != nil {
		t.Fatalf("unmarshal config data: %v", err)
	}

	common, ok := config["common"].(map[string]interface{})
	if !ok {
		t.Fatal("missing or invalid \"common\" key in config")
	}
	for _, field := range []string{"key", "method", "address", "login_password", "auth_password", "log_file", "log_level"} {
		if _, ok := common[field]; !ok {
			t.Errorf("common missing field %q", field)
		}
	}
	if common["key"] != "k1" {
		t.Errorf("expected common.key=k1, got %v", common["key"])
	}
	if common["method"] != "aes-256-cfb" {
		t.Errorf("expected common.method=aes-256-cfb, got %v", common["method"])
	}

	svcSection, ok := config["test-svc"].(map[string]interface{})
	if !ok {
		t.Fatal("missing service-named key \"test-svc\" in config")
	}
	if _, ok := svcSection["key"]; !ok {
		t.Error("service section missing \"key\" field")
	}
}
