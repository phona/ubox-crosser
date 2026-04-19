//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"
)

type versionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
}

func TestVersionEndpoint_ConcurrentRequests(t *testing.T) {
	const numRequests = 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := httpClient().Get(versionURL())
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				results <- fmt.Errorf("status %d", resp.StatusCode)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				results <- err
				return
			}

			var info versionInfo
			if err := json.Unmarshal(body, &info); err != nil {
				results <- err
				return
			}
			results <- nil
		}()
	}

	for i := 0; i < numRequests; i++ {
		if err := <-results; err != nil {
			t.Errorf("concurrent request %d failed: %v", i, err)
		}
	}
}

func TestVersionEndpoint_Idempotent(t *testing.T) {
	var first, second versionInfo

	for i, target := range []*versionInfo{&first, &second} {
		resp, err := httpClient().Get(versionURL())
		if err != nil {
			t.Fatalf("request %d failed: %v", i+1, err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err := json.Unmarshal(body, target); err != nil {
			t.Fatalf("request %d: invalid JSON: %v", i+1, err)
		}
	}

	if first.Version != second.Version || first.Commit != second.Commit || first.BuildTime != second.BuildTime {
		t.Errorf("consecutive requests returned different results:\n  first:  %+v\n  second: %+v", first, second)
	}
}
