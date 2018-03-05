package p9p

import "errors"

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

func negotiateClient(c *Conn, tv *Tversion) (*Rversion, error) {
	tv.size = uint32(4 + 1 + 2 + 4 + len(tv.version.data))

	logf("open: sending %#v\n", tv)
	if err := tv.WriteBinary(c); err != nil {
		logf("open: err: %s\n", err)
		return nil, err
	}

	if err := c.Flush(); err != nil {
		return nil, err
	}

	rv := Rversion{}
	if err := rv.ReadBinary(c); err != nil {
		return nil, err
	}
	logf("open: recv %#v\n", rv)
	return &rv, nil
}

func negotiateServer(c *Conn, tv *Tversion) error {
	rv := Rversion{}
	if err := rv.ReadBinary(c); err != nil {
		return err
	}
	logf("accept: got %#v\n", tv)

	if !supported(string(tv.version.data)) {
		return ErrBadVersion
	}
	if tv.msize != rv.msize {
		tv.msize = min(rv.msize, tv.msize)
	}
	if tv.msize <= 0 {
		return ErrBadIOSize
	}

	tv.size = uint32(4 + 1 + 2 + 4 + len(tv.version.data))
	if err := tv.WriteBinary(c); err != nil {
		logf("open: err: %s\n", err)
		return err
	}
	if err := c.Flush(); err != nil {
		return err
	}
	c.state = StEstablished
	c.version = string(rv.version.data)
	return nil
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
