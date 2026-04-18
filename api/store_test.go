package api

import (
	"sync"
	"testing"

	"github.com/phona/ubox-crosser/models/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				LogFile:  "/var/log/svc1.log",
			},
		},
	}
}

func TestConfigStore_GetAll_ReturnsCopy(t *testing.T) {
	t.Parallel()
	store := NewConfigStore(testConfigs())
	all := store.GetAll()
	all["svc1"] = config.ServerConfig{}
	assert.Equal(t, "debug", store.GetAll()["svc1"].LogLevel, "GetAll must return a copy")
}

func TestConfigStore_Update_MutableFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		field   string
		value   string
		checkFn func(t *testing.T, cfg config.ServerConfig)
	}{
		{"log_level", "log_level", "info", func(t *testing.T, cfg config.ServerConfig) {
			assert.Equal(t, "info", cfg.LogLevel)
		}},
		{"log_file", "log_file", "/tmp/new.log", func(t *testing.T, cfg config.ServerConfig) {
			assert.Equal(t, "/tmp/new.log", cfg.LogFile)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := NewConfigStore(testConfigs())
			updated, err := store.Update(map[string]map[string]string{
				"svc1": {tt.field: tt.value},
			})
			require.NoError(t, err)
			assert.Len(t, updated, 1)
			tt.checkFn(t, store.GetAll()["svc1"])
		})
	}
}

func TestConfigStore_Update_ImmutableFields(t *testing.T) {
	t.Parallel()
	immutable := []string{"key", "method", "address", "login_password", "auth_password"}
	for _, field := range immutable {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			store := NewConfigStore(testConfigs())
			_, err := store.Update(map[string]map[string]string{
				"svc1": {field: "new-value"},
			})
			var immErr *ImmutableFieldError
			require.ErrorAs(t, err, &immErr)
			assert.Equal(t, field, immErr.Field)
		})
	}
}

func TestConfigStore_Update_UnknownService(t *testing.T) {
	t.Parallel()
	store := NewConfigStore(testConfigs())
	_, err := store.Update(map[string]map[string]string{
		"nonexistent": {"log_level": "info"},
	})
	var svcErr *ServiceNotFoundError
	require.ErrorAs(t, err, &svcErr)
}

func TestConfigStore_Update_Atomicity(t *testing.T) {
	t.Parallel()
	store := NewConfigStore(testConfigs())
	_, err := store.Update(map[string]map[string]string{
		"svc1": {"log_level": "info", "key": "bad"},
	})
	require.Error(t, err)
	assert.Equal(t, "debug", store.GetAll()["svc1"].LogLevel, "no field should change on partial failure")
}

func TestConfigStore_ConcurrentAccess(t *testing.T) {
	t.Parallel()
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
			store.Update(map[string]map[string]string{"svc1": {"log_level": "warn"}}) //nolint:errcheck
		}()
	}
	wg.Wait()
}
