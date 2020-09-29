package session

import (
	"bytes"
	"net"
	"context"
)

type Connect interface {
	Close()
	SendMessage(msg string) error

}

// Conn wrap net.Conn
type Conn struct {
	rawConn   net.Conn
	sendCh    chan []byte
	done      chan error
	name      string
	messageCh chan *[]byte
}

// NewConn create new conn
func NewConn(c net.Conn) *Conn {
	conn := &Conn{
		rawConn:   c,
		sendCh:    make(chan []byte),
		done:      make(chan error),
		messageCh: make(chan *[]byte),
	}

	conn.name = c.RemoteAddr().String()

	return conn
}

// Close close connection
func (c *Conn) Close() {
	_ = c.rawConn.Close()
}

func (c *Conn) Write(p []byte) (n int, err error){
	c.sendCh <- p
	return
}

// SendMessage send message
func (c *Conn) SendMessage(buf []byte) error {
	c.sendCh <- buf
	return nil
}

// SendMessage send message
func (c *Conn) SendErr(err error) error {
	c.sendCh <- []byte(err.Error())
	return nil
}

// SendMessage send message
func (c *Conn) SendMsg(msg string) error {
	c.sendCh <- []byte(msg)
	return nil
}

// writeCoroutine write coroutine
func (c *Conn) writeCoroutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case pkt := <-c.sendCh:
			if pkt == nil {
				continue
			}

			if _, err := c.rawConn.Write(pkt); err != nil {
				c.done <- err
			}
		}
	}
}

// readCoroutine read coroutine
func (c *Conn) readCoroutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 读取数据
			buf := make([]byte, 1024)
			n, err := c.rawConn.Read(buf)
			if err != nil {
				c.done <- err
				continue
			}
			if n == 0 {
				continue
			}
			r := bytes.TrimRight(buf, "\x00")
			c.messageCh <- &r
		}
	}
}
