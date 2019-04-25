package conn

import (
	"bufio"
	"io"
	"net"
)

const sep = '\n'

type Coordinator struct {
	Name   string
	Conn   net.Conn
	reader *bufio.Reader
	closed bool
}

// TODO: add cipher in here
func AsCoordinator(conn net.Conn) *Coordinator {
	c := Coordinator{
		Conn:   conn,
		reader: bufio.NewReader(conn),
		closed: false,
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
