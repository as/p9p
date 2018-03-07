package p9p

import (
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
)

var (
	ErrWrongMsg = errors.New("wrong message")
)

type Header struct {
	Size uint32
	Kind byte
}

func (c *Conn) writeMsg(kind byte, tag uint16, p []byte) bool {
	defer c.Flush() // todo
	return c.writeHeader(kind, tag, p) && c.write(p)
}

func (c *Conn) readMsg(kind byte) bool {
	h := Header{Kind: kind}
	return c.readHeader(&h) && c.read(h.Size)
}

func (c *Conn) writeHeader(kind byte, tag uint16, p []byte) bool {
	if c.err != nil {
		return false
	}

	return c.writebinary(struct {
		Size uint32
		Kind byte
		Tag  uint16
	}{4 + 1 + 2, kind, uint16(tag)})
}

func (c *Conn) readHeader(hdr *Header) bool {
	want := hdr.Kind
	if !c.readbinary(hdr) {
		return false
	}

	if want != 0 && want != hdr.Kind {
		c.err = ErrWrongMsg
		return false
	}

	return true
}

func (c *Conn) readbinary(v interface{}) bool {
	if c.err != nil {
		return false
	}
	c.err = binary.Read(c, binary.BigEndian, v)
	return c.err == nil
}

func (c *Conn) writebinary(v interface{}) bool {
	if c.err != nil {
		return false
	}
	c.err = binary.Write(c, binary.BigEndian, v)
	return c.err == nil
}

func (c *Conn) read(size uint32) bool {
	if c.err != nil {
		return false
	}

	c.tmp, c.err = ioutil.ReadAll(io.LimitReader(c, int64(size)))
	return c.err == nil
}

func (c *Conn) write(p []byte) bool {
	if c.err != nil {
		return false
	}

	_, err := c.Write(p)
	return err == nil
}

func (c *Conn) String() string {
	if len(c.tmp) < 4 {
		return ""
	}
	return string(c.tmp)
}

func (c *Conn) readstring() bool {
	var n uint32
	return c.readbinary(&n) && c.read(n)
}

func (c *Conn) writestring(s string) bool {
	return c.writebinary(uint32(len(s))) && c.write([]byte(s))
}
