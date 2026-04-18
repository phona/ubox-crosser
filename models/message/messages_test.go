package message

import (
	"encoding/json"
	"testing"
)

func TestMessage_UnmarshalJSON_ValidLogin(t *testing.T) {
	raw := `{"type":0,"serve_name":"test_svc","password":"secret"}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != LOGIN {
		t.Errorf("Type = %d, want %d", msg.Type, LOGIN)
	}
	if msg.ServeName != "test_svc" {
		t.Errorf("ServeName = %q, want %q", msg.ServeName, "test_svc")
	}
	if msg.Password != "secret" {
		t.Errorf("Password = %q, want %q", msg.Password, "secret")
	}
}

func TestMessage_UnmarshalJSON_UnknownType(t *testing.T) {
	raw := `{"type":99,"serve_name":"svc","password":"pw"}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err == nil {
		t.Error("expected error for unknown type, got nil")
	}
	if msg.Type == 99 {
		t.Error("Type should not be 99 after validation")
	}
}

func TestMessage_UnmarshalJSON_EmptyServeName_Login(t *testing.T) {
	raw := `{"type":0,"serve_name":"","password":"pw"}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err == nil {
		t.Error("expected error for empty serve_name on LOGIN, got nil")
	}
}

func TestMessage_UnmarshalJSON_EmptyPassword_Login(t *testing.T) {
	raw := `{"type":0,"serve_name":"svc","password":""}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err == nil {
		t.Error("expected error for empty password on LOGIN, got nil")
	}
}

func TestMessage_UnmarshalJSON_Heartbeat_NoValidation(t *testing.T) {
	raw := `{"type":1,"serve_name":"","password":""}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Errorf("heartbeat should not require fields: %v", err)
	}
	if msg.Type != HEART_BEAT {
		t.Errorf("Type = %d, want %d", msg.Type, HEART_BEAT)
	}
}

func TestMessage_UnmarshalJSON_GenWorker_NoPasswordRequired(t *testing.T) {
	raw := `{"type":2,"serve_name":"svc","password":""}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Errorf("GEN_WORKER should not require password: %v", err)
	}
}

func TestMessage_UnmarshalJSON_GenWorker_EmptyServeName(t *testing.T) {
	raw := `{"type":2,"serve_name":"","password":""}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err == nil {
		t.Error("expected error for empty serve_name on GEN_WORKER, got nil")
	}
}

func TestResultMessage_UnmarshalJSON_ValidSuccess(t *testing.T) {
	raw := `{"result":4,"reason":0}`
	var msg ResultMessage
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Result != SUCCESS {
		t.Errorf("Result = %d, want %d", msg.Result, SUCCESS)
	}
}

func TestResultMessage_UnmarshalJSON_ValidFailed(t *testing.T) {
	raw := `{"result":5,"reason":1}`
	var msg ResultMessage
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Result != FAILED {
		t.Errorf("Result = %d, want %d", msg.Result, FAILED)
	}
}

func TestResultMessage_UnmarshalJSON_InvalidResultCode(t *testing.T) {
	raw := `{"result":99,"reason":0}`
	var msg ResultMessage
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Fatalf("unmarshal should not return error: %v", err)
	}
	if msg.Result == 99 {
		t.Error("Result should be clamped for unknown code")
	}
}
