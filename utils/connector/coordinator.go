package connector

import (
	"bufio"
	"io"
	"net"
	"sync"
)

const sep = '\n'

type Coordinator struct {
	Name       string
	Conn       net.Conn
	reader     *bufio.Reader
	closed     bool
	readMutex  sync.Mutex
	writeMutex sync.Mutex
}

// TODO: add cipher in here
func AsCoordinator(conn net.Conn) *Coordinator {
	return &Coordinator{
		Conn:   conn,
		reader: bufio.NewReader(conn),
		closed: false,
	}
}

func (c *Coordinator) ReadMsg() (string, error) {
	c.readMutex.Lock()
	defer c.readMutex.Unlock()
	result, err := c.reader.ReadString(sep)
	if err == io.EOF {
		c.closed = true
	}
	return result, err
}

func (c *Coordinator) SendMsg(content string) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	contentBuf := []byte(content)
	contentBuf = append(contentBuf, sep)
	_, err := c.Conn.Write(contentBuf)
	return err
}

func (c *Coordinator) IsTerminate() bool {
	return c.closed
}

func (c *Coordinator) Close() {
	c.closed = true
	c.Conn.Close()
}
