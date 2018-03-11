package p9p

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type ProtocolError struct {
	msg string
}

func (e ProtocolError) Error() string {
	return fmt.Sprintf("9perror: %q", e.msg)
}

var (
	ErrWrongMsg = errors.New("wrong message")
)

type Header struct {
	Size uint32
	Kind
	Tag uint16
}

type Msg struct {
	Header
	src io.Reader
	bytes.Buffer
	err error
}

func (c *Msg) writeMsg(kind Kind, tag uint16, p []byte) bool {
	return c.writeHeader(kind, tag) && c.write(p)
}

func (c *Msg) readMsg(kind Kind) bool {
	return c.readHeader() && c.read(c.Header.Size-4-1-2)
}

func (c *Msg) writeHeader(kind Kind, tag uint16) bool {
	if c.err != nil {
		return false
	}
	c.Header.Tag = tag
	c.Header.Kind = kind

	return c.writebinary(struct {
		Kind Kind
		Tag  uint16
	}{kind, uint16(tag)})
}

func (c *Msg) self() bool {
	c.src = bytes.NewReader(c.Buffer.Bytes())
	c.Buffer = bytes.Buffer{}
	return true
}

func (c *Msg) Transmit(w io.Writer) error {
	c.err = binary.Write(w, binary.LittleEndian, 4+uint32(c.Buffer.Len()))
	if c.err != nil {
		return c.err
	}
	_, c.err = c.Buffer.WriteTo(w)
	if c.err != nil {
		return c.err
	}
	type Flusher interface {
		Flush() error
	}
	if f, ok := w.(Flusher); ok {
		c.err = f.Flush()
	}
	return c.err
}

func (c *Msg) readHeader() bool {
	c.Buffer.Reset()
	if !c.readbinary(&c.Header) {
		return false
	}

	if c.Header.Kind == KRerror {
		c.err = ProtocolError{}
		c.readstring()
		c.err = fmt.Errorf("remote: %q", c.Buffer.Bytes())
		return false
	}

	return true
}

func (c *Msg) readbinary(v interface{}) bool {
	if c.err != nil {
		return false
	}
	c.err = binary.Read(c.src, binary.LittleEndian, v)
	return c.err == nil
}

func (c *Msg) writebinary(v interface{}) bool {
	if c.err != nil {
		return false
	}
	c.err = binary.Write(&c.Buffer, binary.LittleEndian, v)
	return c.err == nil
}

func (c *Msg) read(size uint32) bool {
	if c.err != nil {
		return false
	}
	_, c.err = io.Copy(&(c.Buffer), io.LimitReader(c.src, int64(size)))
	return c.err == nil
}

func (c *Msg) write(p []byte) bool {
	if c.err != nil {
		return false
	}

	_, err := c.Buffer.Write(p)
	return err == nil
}

func (c *Msg) String() string {
	if c.Buffer.Len() < 4 {
		return ""
	}
	return c.Buffer.String()
}

func (c *Msg) readstring() bool {
	var n uint16
	return c.readbinary(&n) && c.read(uint32(n))
}

func (c *Msg) writestring(s string) bool {
	return c.writebinary(uint16(len(s))) && c.write([]byte(s))
}
