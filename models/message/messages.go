package message

import "ubox-crosser/models/errors"

const (
	LOGIN = iota
	HEART_BEAT
	GEN_WORKER
	AUTHENTICATION

	SUCCESS
)

type Message struct {
	Type       uint8  `json:"type"`
	ServeName  string `json:"serve_name"`
	Password   string `json:"password"`
	ServerAddr string `json:"server_addr"`
	Mode       string `json:"mode"`
	// connect address
}

type ResultMessage struct {
	Result errors.Error `json:"result"`
}

type CheckModeRequest struct {
	Password string `json:"password"`
}

type CheckModeResponse struct {
	Result uint8  `json:"result"`
	Reason string `json:"reason"`
}
