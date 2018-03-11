package p9p

import (
	"errors"
)

var (
	ErrBadIOSize = errors.New("bad io size")
)

type entry struct {
	tag    uint16
	expect Kind
}
type msg struct {
	*Msg
	reply chan *Msg
}

func (c *Conn) run() {
	inflight := make(map[uint16]chan *Msg)
	ctr := uint16(1)
	incoming := make(chan Msg)
	go func(){
		for {
			m := Msg{src: c}
			if !m.readMsg(0) {
				logf("readmsg: %s\n", m.err)
				panic("TODO(as): handle readMsg failure")
			}
			logf("readmsg loop: %#v\n", m)
			incoming <- m
		}
	}()
	for{
	select {
	case m := <-c.txout:
		if err := m.Transmit(c); err != nil {
			panic("TODO(as): handle transmit failure")
		}
	logf("txoutdone: %#v\n", m.Msg)
		inflight[m.Tag] = m.reply
		ctr++
	case m := <-incoming:
		repl, ok := inflight[m.Tag]
		if !ok {
			logf("-incoming: %s\n", m.err)
			panic("TODO(as): handle invalid tag sent as reply")
		}
		repl <- &m
	}
	}
}

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

func (c *Conn) schedule(m *Msg) (ok bool) {
	logf("schedule: %#v\n", m)
	defer func(){logf("de schedule: %#v\n", m)}()
	if m.err != nil{
		return false
	}
	reply := make(chan *Msg)
	logf("c.txout<-msg{m, reply}\n")
	c.txout<-msg{m, reply}
	logf("done! c.txout<-msg{m, reply}\n")
	m =  <-reply
	logf("done!! c.txout<-msg{m, reply}\n")
	if m.err == nil{
		m.self()
	}
	return m.err == nil
}

func (c *Conn) Ver() (max uint32, version string, err error) {
	m := &Msg{src: c}
	if !m.writeHeader(KTversion, 0xffff) || !m.writebinary(uint32(0xffff)) || !m.writestring("9P2000") {
		return max, version, m.err
	}
	
	if !c.schedule(m){
		return max, version, m.err
	}
	
	m.readbinary(&max)
	m.readstring()

	// TODO(as): negotiate version by comparing client and server, choose smallest value
	c.version = m.String()
	c.state = StEstablished
	return max, c.version, m.err

}

func (c *Conn) Attach(fid int, afid int, uname, aname string) error {
	m := &Msg{src: c}
	m.writeHeader(KTattach, 1)
	m.writebinary(uint16(fid))
	m.writebinary(uint16(afid))
	m.writestring(uname)
	m.writestring(aname)

	logf("attach: %v %s %s", m.Header, m.String(), m.err)
	if !c.schedule(m){
		return m.err
	}
	
	logf("attach: %v %s %s", m.Header, m.String(), m.err)
	m.read(13)
	logf("attach: %v %s %s", m.Header, m.String(), m.err)
	return m.err
}

func negotiateClient(c *Conn, tv *Tversion) (*Rversion, error) {
	tv.size = uint32(4 + 1 + 2 + 4 + 2 + len(tv.version.data))

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

	tv.size = uint32(4 + 1 + 2 + 4 + 2 + len(tv.version.data))
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
