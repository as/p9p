package p9p

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

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

type Names struct {
	List []string
}

func (m *Msg) String() string {
	return fmt.Sprintf("size[%d] kind[%s] tag[%d] -> data[%q]\n",
		m.Header.Size,
		m.Header.Kind,
		m.Header.Tag,
		m.Buffer.Bytes(),
	)
}

func (c *Msg) size() uint32 {
	return 4 + 1 + 2 + uint32(c.Buffer.Len())
}

func (c *Msg) write(p []byte) bool {
	if c.err == nil {
		_, c.err = c.Buffer.Write(p)
	}
	return c.err == nil
}

func (c *Msg) read(size uint32) bool {
	if c.err == nil {
		_, c.err = io.Copy(&(c.Buffer), io.LimitReader(c.src, int64(size)))
	}
	return c.err == nil
}

func (c *Msg) writebinary(v interface{}) bool {
	if c.err == nil {
		c.err = binary.Write(&c.Buffer, binary.LittleEndian, v)
	}
	return c.err == nil
}

func (c *Msg) readbinary(v interface{}) bool {
	if c.err == nil {
		c.err = binary.Read(c.src, binary.LittleEndian, v)
	}
	return c.err == nil
}

func (c *Msg) writeMsg(kind Kind, p []byte) bool {
	return c.writeHeader(kind) && c.write(p)
}

func (c *Msg) readMsg(kind Kind) bool {
	return c.readHeader() && c.read(c.Header.Size-4-1-2)
}

func (c *Msg) writeHeader(kind Kind) bool {
	if c.err == nil {
		c.Header.Kind = kind
	}
	return c.err == nil
}

func (c *Msg) readHeader() bool {
	if c.err == nil {
		c.Buffer.Reset()
	}
	return c.readbinary(&c.Header)
}

func (c *Msg) readbytes(p []byte) bool {
	if c.err == nil {
		_, c.err = io.ReadAtLeast(&c.Buffer, p, len(p))
	}
	return c.err == nil
}

func (c *Msg) writestring(s string) bool {
	return c.writebinary(uint16(len(s))) && c.write([]byte(s))
}

func (c *Msg) readstring() bool {
	var n uint16
	return c.readbinary(&n) && c.read(uint32(n))
}

func (m *Msg) readQuids(q *[]Qid) bool {
	var (
		nn uint16
		q0 Qid
	)
	m.readbinary(&nn)
	for i := uint16(0); i < nn; i++ {
		if !m.readbinary(&q0) {
			break
		}
		*q = append(*q, q0)
	}
	return m.err == nil

}

func (m *Msg) writeNames(names ...string) bool {
	m.writebinary(uint16(len(names)))
	for i := range names {
		if !m.writestring(names[i]) {
			break
		}
	}
	return m.err == nil
}

func (m *Msg) readNames(nm *Names) bool {
	var (
		nn, n uint16
		b     [65536]byte
	)
	m.readbinary(&nn)
	for i := uint16(0); i < nn; i++ {
		if !m.readbinary(&n) || !m.readbytes(b[:n]) {
			break
		}
		nm.List = append(nm.List, string(b[:n]))
	}
	return m.err == nil
}

func (c *Msg) self() bool {
	c.src = bytes.NewReader(c.Buffer.Bytes())
	c.Buffer = bytes.Buffer{}
	return c.err == nil
}
