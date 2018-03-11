package p9p

import (
	"bufio"
	"encoding/binary"
	"net"
)

type Config struct {
	// IOUnit is the maximum size of a 9p message, the protocol negotiates
	// this during exchange of the version message and selects the minimum
	// supported size between client and server
	IOUnit int

	// Protocol strings that match the regexp 9P[0-9]+
	// the default value is 9P2000. Older versions are not supported
	// at this time
	Protocol string
}

const (
	Version = "9P2000"
	MaxMsg  = 65535
)

var defaultConfig = Config{
	IOUnit:   65536,
	Protocol: "9P2000",
}

type Conn struct {
	txout   chan msg
	netconn net.Conn
	bio     *bufio.ReadWriter
	err     error
	version string
	state   State
	exit    chan bool
	done    chan struct{}
}

// NewConn opens a new 9p connection from an existing net.Conn.
func NewConn(conn net.Conn, conf *Config) (c *Conn) {
	if conf == nil {
		conf = &defaultConfig
	}
	defer func() {
		c.exit <- true
		go c.run() // see run.go:/run/
	}()
	return &Conn{
		txout:   make(chan msg),
		netconn: conn,
		exit:    make(chan bool, 1),
		done:    make(chan struct{}),
		bio: bufio.NewReadWriter(
			bufio.NewReaderSize(conn, conf.IOUnit),
			bufio.NewWriterSize(conn, conf.IOUnit),
		),
	}
}

func (c *Conn) Transmit(m *Msg) (err error) {
	defer func() { logf("called Transmit, err=%s", err) }()
	m.Header.Size = m.size()

	if err = binary.Write(c, binary.LittleEndian, m.Header); err != nil {
		return c.err
	}

	if _, err = m.Buffer.WriteTo(c); err != nil {
		return c.err
	}

	return c.bio.Flush()
}

func (c *Conn) Read(p []byte) (n int, err error) {
	defer func() { logf("called conn.Read, result n=%d, err=%v", n, err) }()
	return c.bio.Read(p)
}

func (c *Conn) Write(p []byte) (n int, err error) {
	defer func() { logf("called conn.Write: %q\n\tresult n=%d, err=%v", p, n, err) }()
	return c.bio.Write(p)
}

func (c *Conn) close(cleanup bool) (err error) {
	if cleanup {
		close(c.exit)
		close(c.done)
		return c.netconn.Close()
	}
	return nil
}

func (c *Conn) Close() (err error) {
	defer func() { logf("called conn.Close, result err=%s", err) }()
	return c.close(<-c.exit)
}
