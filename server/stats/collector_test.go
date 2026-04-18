package stats

import (
	"sync"
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	c := NewCollector()
	snap := c.Snapshot()
	if snap.TotalConnections != 0 {
		t.Errorf("expected TotalConnections=0, got %d", snap.TotalConnections)
	}
	if snap.ActiveConnections != 0 {
		t.Errorf("expected ActiveConnections=0, got %d", snap.ActiveConnections)
	}
	if snap.TotalBytesIn != 0 {
		t.Errorf("expected TotalBytesIn=0, got %d", snap.TotalBytesIn)
	}
	if snap.TotalBytesOut != 0 {
		t.Errorf("expected TotalBytesOut=0, got %d", snap.TotalBytesOut)
	}
	if snap.AvgTunnelDurationMs != 0 {
		t.Errorf("expected AvgTunnelDurationMs=0, got %d", snap.AvgTunnelDurationMs)
	}
}

func TestTrackAndReleaseConnection(t *testing.T) {
	c := NewCollector()

	c.TrackConnection()
	snap := c.Snapshot()
	if snap.TotalConnections != 1 {
		t.Errorf("expected TotalConnections=1, got %d", snap.TotalConnections)
	}
	if snap.ActiveConnections != 1 {
		t.Errorf("expected ActiveConnections=1, got %d", snap.ActiveConnections)
	}

	c.TrackConnection()
	snap = c.Snapshot()
	if snap.TotalConnections != 2 {
		t.Errorf("expected TotalConnections=2, got %d", snap.TotalConnections)
	}
	if snap.ActiveConnections != 2 {
		t.Errorf("expected ActiveConnections=2, got %d", snap.ActiveConnections)
	}

	startTime := time.Now().Add(-100 * time.Millisecond)
	c.ReleaseConnection(startTime)
	snap = c.Snapshot()
	if snap.ActiveConnections != 1 {
		t.Errorf("expected ActiveConnections=1, got %d", snap.ActiveConnections)
	}
	if snap.TotalConnections != 2 {
		t.Errorf("expected TotalConnections=2 (unchanged), got %d", snap.TotalConnections)
	}
	if snap.AvgTunnelDurationMs < 100 {
		t.Errorf("expected AvgTunnelDurationMs >= 100, got %d", snap.AvgTunnelDurationMs)
	}
}

func TestAddBytes(t *testing.T) {
	c := NewCollector()

	c.AddBytesIn(1024)
	c.AddBytesIn(512)
	c.AddBytesOut(2048)

	snap := c.Snapshot()
	if snap.TotalBytesIn != 1536 {
		t.Errorf("expected TotalBytesIn=1536, got %d", snap.TotalBytesIn)
	}
	if snap.TotalBytesOut != 2048 {
		t.Errorf("expected TotalBytesOut=2048, got %d", snap.TotalBytesOut)
	}
}

func TestAvgDurationMultipleConnections(t *testing.T) {
	c := NewCollector()

	// Simulate two connections with known durations
	c.TrackConnection()
	c.TrackConnection()

	start1 := time.Now().Add(-200 * time.Millisecond)
	start2 := time.Now().Add(-400 * time.Millisecond)

	c.ReleaseConnection(start1)
	c.ReleaseConnection(start2)

	snap := c.Snapshot()
	if snap.ActiveConnections != 0 {
		t.Errorf("expected ActiveConnections=0, got %d", snap.ActiveConnections)
	}
	// Average should be roughly (200+400)/2 = 300ms, allow wide tolerance
	if snap.AvgTunnelDurationMs < 200 {
		t.Errorf("expected AvgTunnelDurationMs >= 200, got %d", snap.AvgTunnelDurationMs)
	}
}

func TestConcurrentAccess(t *testing.T) {
	c := NewCollector()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.TrackConnection()
			c.AddBytesIn(10)
			c.AddBytesOut(20)
			c.ReleaseConnection(time.Now())
		}()
	}
	wg.Wait()

	snap := c.Snapshot()
	if snap.TotalConnections != 100 {
		t.Errorf("expected TotalConnections=100, got %d", snap.TotalConnections)
	}
	if snap.ActiveConnections != 0 {
		t.Errorf("expected ActiveConnections=0, got %d", snap.ActiveConnections)
	}
	if snap.TotalBytesIn != 1000 {
		t.Errorf("expected TotalBytesIn=1000, got %d", snap.TotalBytesIn)
	}
	if snap.TotalBytesOut != 2000 {
		t.Errorf("expected TotalBytesOut=2000, got %d", snap.TotalBytesOut)
	}
}
