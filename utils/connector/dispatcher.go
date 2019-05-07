package connector

import (
	"net"
	"sync"
)

type Dispatcher struct {
	listeners []net.Listener
	errs      chan error
	conns     chan net.Conn

	mutex sync.Mutex
}

func NewDispatcher(size uint64) *Dispatcher {
	return &Dispatcher{
		listeners: make([]net.Listener, 0, size),
		errs:      make(chan error, 10),
		conns:     make(chan net.Conn, 10),
	}
}

func (d *Dispatcher) Add(listener net.Listener) {
	d.mutex.Lock()
	d.listeners = append(d.listeners, listener)
	d.mutex.Unlock()
	go d.listen(listener)
}

func (d *Dispatcher) listen(listener net.Listener) {
	for {
		if conn, err := listener.Accept(); err != nil {
			d.errs <- err
		} else {
			d.conns <- conn
		}
	}
}

func (d *Dispatcher) Err() error {
	return <-d.errs
}

func (d *Dispatcher) Conn() net.Conn {
	return <-d.conns
}
