package utils

import "encoding/json"

func MakeJsonBuf(i interface{}) []byte {
	buf, _ := json.Marshal(i)
	return buf
}
