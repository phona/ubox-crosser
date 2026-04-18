package server

import (
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
		t.Errorf("expected AvgTunnelDurationMs=0, got %f", snap.AvgTunnelDurationMs)
	}
}

func TestCollectorTunnelLifecycle(t *testing.T) {
	c := NewCollector()

	c.OnTunnelStart()
	snap := c.Snapshot()
	if snap.TotalConnections != 1 {
		t.Errorf("expected TotalConnections=1, got %d", snap.TotalConnections)
	}
	if snap.ActiveConnections != 1 {
		t.Errorf("expected ActiveConnections=1, got %d", snap.ActiveConnections)
	}

	start := time.Now().Add(-100 * time.Millisecond)
	c.OnTunnelEnd(start)
	snap = c.Snapshot()
	if snap.TotalConnections != 1 {
		t.Errorf("expected TotalConnections=1, got %d", snap.TotalConnections)
	}
	if snap.ActiveConnections != 0 {
		t.Errorf("expected ActiveConnections=0, got %d", snap.ActiveConnections)
	}
	if snap.AvgTunnelDurationMs < 100 {
		t.Errorf("expected AvgTunnelDurationMs >= 100, got %f", snap.AvgTunnelDurationMs)
	}
}

func TestCollectorMultipleTunnels(t *testing.T) {
	c := NewCollector()

	c.OnTunnelStart()
	c.OnTunnelStart()
	snap := c.Snapshot()
	if snap.ActiveConnections != 2 {
		t.Errorf("expected ActiveConnections=2, got %d", snap.ActiveConnections)
	}

	start := time.Now().Add(-50 * time.Millisecond)
	c.OnTunnelEnd(start)
	snap = c.Snapshot()
	if snap.ActiveConnections != 1 {
		t.Errorf("expected ActiveConnections=1, got %d", snap.ActiveConnections)
	}
	if snap.TotalConnections != 2 {
		t.Errorf("expected TotalConnections=2, got %d", snap.TotalConnections)
	}
}

func TestCollectorByteCounters(t *testing.T) {
	c := NewCollector()

	c.AddBytesIn(1024)
	c.AddBytesIn(2048)
	c.AddBytesOut(512)

	snap := c.Snapshot()
	if snap.TotalBytesIn != 3072 {
		t.Errorf("expected TotalBytesIn=3072, got %d", snap.TotalBytesIn)
	}
	if snap.TotalBytesOut != 512 {
		t.Errorf("expected TotalBytesOut=512, got %d", snap.TotalBytesOut)
	}
}

func TestCollectorAvgDurationMultiple(t *testing.T) {
	c := NewCollector()

	c.OnTunnelStart()
	start1 := time.Now().Add(-200 * time.Millisecond)
	c.OnTunnelEnd(start1)

	c.OnTunnelStart()
	start2 := time.Now().Add(-100 * time.Millisecond)
	c.OnTunnelEnd(start2)

	snap := c.Snapshot()
	// Average should be roughly (200+100)/2 = 150ms
	if snap.AvgTunnelDurationMs < 100 || snap.AvgTunnelDurationMs > 300 {
		t.Errorf("expected AvgTunnelDurationMs ~150, got %f", snap.AvgTunnelDurationMs)
	}
}

func TestCollectorConcurrentAccess(t *testing.T) {
	c := NewCollector()
	done := make(chan struct{})

	for i := 0; i < 100; i++ {
		go func() {
			c.OnTunnelStart()
			c.AddBytesIn(10)
			c.AddBytesOut(5)
			c.OnTunnelEnd(time.Now())
			done <- struct{}{}
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}

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
	if snap.TotalBytesOut != 500 {
		t.Errorf("expected TotalBytesOut=500, got %d", snap.TotalBytesOut)
	}
}
