package p9p

import (
	"bytes"
	"errors"
	"io"
	"time"
)

var (
	ErrShortRead = errors.New("short read")
)

type Conn struct {
	buf     *bytes.Buffer
	fd      io.ReadWriteCloser
	t       *time.Timer
	lastHDR Hdr
}

func Open(fd io.ReadWriteCloser, opts *Opts) (conn *Conn, err error) {
	if opts == nil {
		opts = &defaultOpts
	}
	conn = &Conn{
		buf: new(bytes.Buffer),
		fd:  fd,
	}
	return conn, nil
}

func (c *Conn) ReadMsg() (msg Msg, err error) {
	err = c.readHDR()
	if err != nil {
		return nil, err
	}
	return msg, err
}
func (c *Conn) WriteMsg(msg Msg) (err error) {
	return msg, err
}

func (c *Conn) readN(n int) (err error) {
	c.ensure(n)
	m, err := io.CopyN(c.buf, c.fd, int64(n))
	if int64(n) != m {
		return ErrShortRead
	}
	return err
}

func (c *Conn) readHDR() (err error) {
	const len = 5
	c.ensure(len)
	n, err := io.CopyN(c.buf, c.fd, len)
	if n < len-1 {
		return ErrShortRead
	}
	return (&c.lastHDR).ReadBinary(c.buf)
}

func (c *Conn) ensure(n int) {
	if m := c.buf.Len(); m < n {
		c.buf.Grow(round(n - m))
	}
}

func round(n int) int {
	return (n + 256) &^ (256)
}
