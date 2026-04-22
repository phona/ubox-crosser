package contract

import (
	"encoding/json"
	"net/http"
)

func writeSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

func writeError(w http.ResponseWriter, status int, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    code,
		"message": message,
	})
}

func readJSON(r *http.Request, v interface{}) error {
	defer func() { _ = r.Body.Close() }()
	return json.NewDecoder(r.Body).Decode(v)
}

var stubServices = map[string]Service{}

func resetStubServices() {
	stubServices = map[string]Service{}
}

func stubListServicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, 9999, "method not allowed")
		return
	}
	svcs := make([]Service, 0, len(stubServices))
	for _, s := range stubServices {
		svcs = append(svcs, s)
	}
	writeSuccess(w, http.StatusOK, ServiceListData{Services: svcs, Total: len(svcs)})
}
