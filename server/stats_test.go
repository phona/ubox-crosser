package server

import (
	"testing"
)

func TestNewCollector_ZeroCounters(t *testing.T) {
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

func TestCollector_RecordConnection(t *testing.T) {
	c := NewCollector()
	c.RecordConnection()
	c.RecordConnection()
	c.RecordConnection()

	snap := c.Snapshot()
	if snap.TotalConnections != 3 {
		t.Errorf("expected TotalConnections=3, got %d", snap.TotalConnections)
	}
}

func TestCollector_TunnelLifecycle(t *testing.T) {
	c := NewCollector()

	c.RecordTunnelStart()
	c.RecordTunnelStart()
	snap := c.Snapshot()
	if snap.ActiveConnections != 2 {
		t.Errorf("expected ActiveConnections=2, got %d", snap.ActiveConnections)
	}

	// End one tunnel with 100ms duration
	c.RecordTunnelEnd(100_000_000) // 100ms in nanoseconds
	snap = c.Snapshot()
	if snap.ActiveConnections != 1 {
		t.Errorf("expected ActiveConnections=1, got %d", snap.ActiveConnections)
	}
	if snap.AvgTunnelDurationMs != 100.0 {
		t.Errorf("expected AvgTunnelDurationMs=100, got %f", snap.AvgTunnelDurationMs)
	}

	// End second tunnel with 200ms duration
	c.RecordTunnelEnd(200_000_000) // 200ms
	snap = c.Snapshot()
	if snap.ActiveConnections != 0 {
		t.Errorf("expected ActiveConnections=0, got %d", snap.ActiveConnections)
	}
	// Average of 100ms and 200ms = 150ms
	if snap.AvgTunnelDurationMs != 150.0 {
		t.Errorf("expected AvgTunnelDurationMs=150, got %f", snap.AvgTunnelDurationMs)
	}
}

func TestCollector_RecordBytes(t *testing.T) {
	c := NewCollector()
	c.RecordBytesIn(100)
	c.RecordBytesIn(200)
	c.RecordBytesOut(50)
	c.RecordBytesOut(75)

	snap := c.Snapshot()
	if snap.TotalBytesIn != 300 {
		t.Errorf("expected TotalBytesIn=300, got %d", snap.TotalBytesIn)
	}
	if snap.TotalBytesOut != 125 {
		t.Errorf("expected TotalBytesOut=125, got %d", snap.TotalBytesOut)
	}
}
