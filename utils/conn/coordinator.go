package conn

import (
	"bufio"
	"io"
	"net"
	"time"
)

const sep = '\n'

type Coordinator struct {
	Name    string
	Conn    net.Conn
	reader  *bufio.Reader
	Timeout time.Duration
	closed  bool
}

// TODO: add cipher in here
func AsCoordinator(conn net.Conn, duration time.Duration) *Coordinator {
	c := Coordinator{
		Conn:   conn,
		reader: bufio.NewReader(conn),
		closed: false,
	}
	if duration == 0 {
		c.Timeout = 30 * time.Second
	} else {
		c.Timeout = duration
	}
	return &c
}

func (c *Coordinator) ReadMsg() (string, error) {
	result, err := c.reader.ReadString(sep)
	if err == io.EOF {
		c.closed = true
	}
	return result, err
}

func (c *Coordinator) SendMsg(content string) error {
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
