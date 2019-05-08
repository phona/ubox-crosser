package message

import "ubox-crosser/models/errors"

const (
	LOGIN = iota
	HEART_BEAT
	GEN_WORKER
	AUTHENTICATION

	SUCCESS
	FAILED
)

type Message struct {
	Type      uint8  `json:"type"`
	ServeName string `json:"serve_name"`
	Password  string `json:"password"`
}

type ResultMessage struct {
	Result uint8        `json:"result"`
	Reason errors.Error `json:"reason"`
}
