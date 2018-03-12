package p9p

import (
	"errors"
)

var (
	ErrBadIOSize = errors.New("bad io size")
)

// Version returns the version of the 9p protocol used the established connection. If
// the connection is dead, or has never been established, or there is a version
// disagreement between the client and server, an error is returned
func (c *Conn) Version() (string, error) {
	if c.state != StEstablished {
		return "", ErrNoConn
	}
	if c.version == "" {
		return "", ErrBadVersion
	}
	return c.version, nil
}

func (c *Conn) Ver() (max uint32, version string, err error) {
	defer logf("Ver: %v %s %s", err)
	m := &Msg{src: c}
	if !m.writeHeader(KTversion) || !m.writebinary(uint32(0xffff)) || !m.writestring("9P2000") {
		return max, version, m.err
	}

	if !c.schedule(m) {
		logf("!c.schedule: %s\n", m.err)
		return max, version, m.err
	}

	m.readbinary(&max)
	logf("c.readbinary: %s\n", m.err)
	m.readstring()
	logf("c.readbinary: %s\n", m.err)

	// TODO(as): negotiate version by comparing client and server, choose smallest value
	c.version = m.String()
	c.state = StEstablished
	return max, c.version, m.err
}

func (c *Conn) Attach(fid, afid Fid, uname, aname string) (q Qid, err error) {
	m := &Msg{src: c}
	m.writeHeader(KTattach)
	m.writebinary(fid)
	m.writebinary(afid)
	m.writestring(uname)
	m.writestring(aname)

	if !c.schedule(m) {
		return q, m.err
	}

	logf("attach: %v %s %s", m.Header, m.String(), m.err)
	m.readbinary(&q)
	logf("attach: %v %s %s", m.Header, m.String(), m.err)
	return q, m.err
}

func supported(ver string) bool {
	return ver == "9p2000"
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
