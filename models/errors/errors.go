package errors

const (
	OK = iota
	INVALID_PASSWORD
	INVALID_SERVE_NAME
	UNKNOWN_CODE
)

func init() {
	customErrors[INVALID_PASSWORD] = "invalid password"
	customErrors[INVALID_SERVE_NAME] = "invalid serve name"
	customErrors[UNKNOWN_CODE] = "unknown code"
}

var customErrors = map[uint8]string{}

type Error uint8

func DecodeError(code uint8) Error {
	return Error(code)
}

func (l Error) Error() string {
	return customErrors[uint8(l)]
}
