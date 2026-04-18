package api

import (
	"sync"
	"testing"

	"github.com/phona/ubox-crosser/models/config"
	"github.com/stretchr/testify/assert"
)

func testConfigs() map[string]config.ServerConfig {
	return map[string]config.ServerConfig{
		"svc1": {
			LoginPass: "pass1",
			AuthPass:  "auth1",
			Address:   "0.0.0.0:7000",
			Config: config.Config{
				Key:      "key1",
				Method:   "chacha20",
				LogLevel: "debug",
				LogFile:  "/var/log/test.log",
			},
		},
	}
}

func TestNewConfigStore_CopiesMap(t *testing.T) {
	orig := testConfigs()
	store := NewConfigStore(orig)
	orig["svc1"] = config.ServerConfig{Address: "changed"}
	got := store.GetAll()
	assert.Equal(t, "0.0.0.0:7000", got["svc1"].Address, "store should not be affected by external map mutation")
}

func TestConfigStore_GetAll_ReturnsCopy(t *testing.T) {
	store := NewConfigStore(testConfigs())
	a := store.GetAll()
	a["svc1"] = config.ServerConfig{Address: "mutated"}
	b := store.GetAll()
	assert.Equal(t, "0.0.0.0:7000", b["svc1"].Address, "GetAll should return independent copies")
}

func TestConfigStore_HasService(t *testing.T) {
	store := NewConfigStore(testConfigs())
	assert.True(t, store.HasService("svc1"))
	assert.False(t, store.HasService("nonexistent"))
}

func TestConfigStore_Update_MutableFields(t *testing.T) {
	store := NewConfigStore(testConfigs())
	store.Update("svc1", map[string]string{"log_level": "warn", "log_file": "/tmp/new.log"})
	got := store.GetAll()
	assert.Equal(t, "warn", got["svc1"].LogLevel)
	assert.Equal(t, "/tmp/new.log", got["svc1"].LogFile)
	// Immutable fields unchanged.
	assert.Equal(t, "key1", got["svc1"].Key)
}

func TestIsMutable(t *testing.T) {
	assert.True(t, IsMutable("log_level"))
	assert.True(t, IsMutable("log_file"))
	assert.False(t, IsMutable("key"))
	assert.False(t, IsMutable("method"))
	assert.False(t, IsMutable("address"))
	assert.False(t, IsMutable("login_password"))
	assert.False(t, IsMutable("auth_password"))
}

func TestConfigStore_ConcurrentReadWrite(t *testing.T) {
	store := NewConfigStore(testConfigs())
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			store.GetAll()
		}()
		go func() {
			defer wg.Done()
			store.Update("svc1", map[string]string{"log_level": "info"})
		}()
	}
	wg.Wait()
	got := store.GetAll()
	assert.Equal(t, "info", got["svc1"].LogLevel)
}
