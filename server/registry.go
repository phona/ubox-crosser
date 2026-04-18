package server

import (
	"sync"
	"time"
	"ubox-crosser/api"
	"ubox-crosser/models/connection"
)

type ConnectionRegistry struct {
	mu    sync.RWMutex
	conns map[string]*connection.Connection
}

func NewConnectionRegistry() *ConnectionRegistry {
	return &ConnectionRegistry{
		conns: make(map[string]*connection.Connection),
	}
}

func (r *ConnectionRegistry) Add(conn *connection.Connection) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conns[conn.ID] = conn
}

func (r *ConnectionRegistry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.conns, id)
}

func (r *ConnectionRegistry) Get(id string) (*connection.Connection, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	conn, ok := r.conns[id]
	return conn, ok
}

func (r *ConnectionRegistry) List(filter api.FilterOptions) ([]*connection.Connection, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*connection.Connection
	for _, conn := range r.conns {
		if filter.ServeName != "" && conn.ServeName != filter.ServeName {
			continue
		}
		if filter.Status != "" && conn.Status != filter.Status {
			continue
		}
		if filter.Type != "" && conn.Type != filter.Type {
			continue
		}
		filtered = append(filtered, conn)
	}

	total := len(filtered)

	// Paginate
	start := (filter.Page - 1) * filter.PageSize
	if start >= total {
		return []*connection.Connection{}, total
	}
	end := start + filter.PageSize
	if end > total {
		end = total
	}

	return filtered[start:end], total
}

func (r *ConnectionRegistry) UpdateHeartbeat(id string, t time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if conn, ok := r.conns[id]; ok {
		conn.LastHeartbeat = &t
	}
}

func (r *ConnectionRegistry) UpdateStatus(id string, status connection.Status) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if conn, ok := r.conns[id]; ok {
		conn.Status = status
	}
}
