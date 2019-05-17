package protocols

import (
	"bufio"
	"io"
	"net"
	"sync"
)

const sep = '\n'

type LineBaseProtocol struct {
	conn       net.Conn
	reader     *bufio.Reader
	closed     bool
	readMutex  sync.Mutex
	writeMutex sync.Mutex
}

// TODO: add cipher in here
func NewLineBaseProtocol(conn net.Conn) *LineBaseProtocol {
	p := MakeLineBaseProtocol(conn)
	return &p
}

func MakeLineBaseProtocol(conn net.Conn) LineBaseProtocol {
	return LineBaseProtocol{
		conn:   conn,
		reader: bufio.NewReader(conn),
		closed: false,
	}
}

func (c *LineBaseProtocol) ReadMsg() ([]byte, error) {
	c.readMutex.Lock()
	defer c.readMutex.Unlock()
	result, err := c.reader.ReadBytes(sep)
	if err == io.EOF {
		c.closed = true
	}
	return result, err
}

func (c *LineBaseProtocol) SendMsg(content []byte) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	contentBuf := content
	contentBuf = append(contentBuf, sep)
	_, err := c.conn.Write(contentBuf)
	return err
}

func (c *LineBaseProtocol) IsTerminate() bool {
	return c.closed
}

func (c *LineBaseProtocol) Close() {
	c.closed = true
	c.conn.Close()
}
