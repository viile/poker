package server

import (
	"bytes"
	"net"
	"context"
)

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

// SendMessage send message
func (c *Conn) SendMessage(buf []byte) error {
	c.sendCh <- buf
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
