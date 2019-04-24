package models

const (
	LOGIN = iota
	HEART_BEAT
	GEN_WORKER
)

type Message struct {
	Type int64  `json:"type"`
	Msg  string `json:"msg"`
	Name string `json:"name"`
}
