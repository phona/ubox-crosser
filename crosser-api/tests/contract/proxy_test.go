package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var stubProxyInstances = map[string]ProxyInstance{}
var stubProxyTokens = map[string]string{}

func resetStubProxy() {
	stubProxyInstances = map[string]ProxyInstance{}
	stubProxyTokens = map[string]string{}
}

func stubProxyMux() *http.ServeMux {
	resetStubProxy()
	resetStubServices()
	stubServices["test-svc"] = Service{
		Name: "test-svc", Key: "k1", Method: "aes-256-cfb", Address: ":8388",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/proxy/register", stubProxyTokenMiddleware(stubRegisterHandler))
	mux.HandleFunc("/api/v1/proxy/heartbeat", stubProxyTokenMiddleware(stubHeartbeatHandler))
	return mux
}

func stubProxyTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Proxy-Token")
		if token == "" {
			writeError(w, http.StatusUnauthorized, 3001, "missing proxy token")
			return
		}
		if token != "valid-service-token" && stubProxyTokens[token] == "" {
			writeError(w, http.StatusUnauthorized, 3001, "invalid proxy token")
			return
		}
		next(w, r)
	}
}

func stubRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, 9999, "method not allowed")
		return
	}
	var req ProxyRegisterRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, 2003, "invalid body")
		return
	}
	token := "proxy-token-" + req.InstanceID
	stubProxyTokens[token] = req.InstanceID
	stubProxyInstances[req.InstanceID] = ProxyInstance{
		InstanceID:  req.InstanceID,
		ServiceName: req.ServiceName,
		Address:     req.Address,
		Status:      "online",
	}
	writeSuccess(w, http.StatusOK, ProxyRegisterData{Token: token})
}

func stubHeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, 9999, "method not allowed")
		return
	}
	var req ProxyHeartbeatRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, 2003, "invalid body")
		return
	}
	if inst, ok := stubProxyInstances[req.InstanceID]; ok {
		inst.ConnectionCount = req.ConnectionCount
		inst.Status = "online"
		stubProxyInstances[req.InstanceID] = inst
	}
	writeSuccess(w, http.StatusOK, nil)
}

// REQ-997-S12: POST /api/v1/proxy/register returns a proxy token
func TestProxyRegister(t *testing.T) {
	srv := httptest.NewServer(stubProxyMux())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodPost, "/api/v1/proxy/register",
		jsonBody(t, ProxyRegisterRequest{
			ServiceName: "test-svc",
			InstanceID:  "inst-01",
			Address:     "10.0.0.1:8388",
		}),
		map[string]string{"X-Proxy-Token": "valid-service-token"})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 0 {
		t.Fatalf("expected code 0, got %d", ar.Code)
	}

	data := parseData[ProxyRegisterData](t, ar)
	if data.Token == "" {
		t.Error("expected non-empty token")
	}
}

// REQ-997-S13: POST /api/v1/proxy/heartbeat updates proxy status
func TestProxyHeartbeat(t *testing.T) {
	srv := httptest.NewServer(stubProxyMux())
	defer srv.Close()

	regResp := doRequest(t, srv, http.MethodPost, "/api/v1/proxy/register",
		jsonBody(t, ProxyRegisterRequest{
			ServiceName: "test-svc",
			InstanceID:  "inst-01",
			Address:     "10.0.0.1:8388",
		}),
		map[string]string{"X-Proxy-Token": "valid-service-token"})
	regAr := parseResponse(t, regResp)
	regData := parseData[ProxyRegisterData](t, regAr)

	resp := doRequest(t, srv, http.MethodPost, "/api/v1/proxy/heartbeat",
		jsonBody(t, ProxyHeartbeatRequest{
			InstanceID:      "inst-01",
			ConnectionCount: 5,
		}),
		map[string]string{"X-Proxy-Token": regData.Token})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 0 {
		t.Fatalf("expected code 0, got %d", ar.Code)
	}
}

// REQ-997-S14: Proxy endpoints reject requests without valid X-Proxy-Token
func TestProxyRejectsNoToken(t *testing.T) {
	srv := httptest.NewServer(stubProxyMux())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodPost, "/api/v1/proxy/register",
		jsonBody(t, ProxyRegisterRequest{
			ServiceName: "test-svc",
			InstanceID:  "inst-01",
			Address:     "10.0.0.1:8388",
		}), nil)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 3001 {
		t.Errorf("expected code 3001, got %d", ar.Code)
	}
}

func TestProxyRejectsInvalidToken(t *testing.T) {
	srv := httptest.NewServer(stubProxyMux())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodPost, "/api/v1/proxy/heartbeat",
		jsonBody(t, ProxyHeartbeatRequest{InstanceID: "inst-01", ConnectionCount: 0}),
		map[string]string{"X-Proxy-Token": "bogus-token"})

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 3001 {
		t.Errorf("expected code 3001, got %d", ar.Code)
	}
}
