package server

import (
	"testing"
	"ubox-crosser/models/config"
)

func TestCCPServer_Run(t *testing.T) {
	ccpConfig := config.CCPConfig{
		Address: "127.0.0.1:7000",
	}
	server := NewCCPServer(ccpConfig)
	if err := server.Run(); err != nil {
		t.Error(err)
		t.Fatal()
	}
}
