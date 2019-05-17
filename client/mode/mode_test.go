package mode

import "testing"

func TestGetConnectMode(t *testing.T) {
	if mode, err := GetConnectMode("127.0.0.1:7000", "", nil); err != nil {
		t.Error(err)
		t.Fatal()
	} else {
		t.Log(mode.Mode())
	}
}
