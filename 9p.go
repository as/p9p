package p9p

import (
	"bufio"
	"errors"
	"net"
)

const (
	Version = "9P2000"
	MaxMsg  = 65535
)

type Kind byte

const (
	KTversion, KRversion Kind = iota + 100, iota + 101
	KTauth, KRauth
	KTattach, KRattach
	KTerror, KRerror
	KTflush, KRflush
	KTwalk, KRwalk
	KTopen, KRopen
	KTcreate, KRcreate
	KTread, KRread
	KTwrite, KRwrite
	KTclunk, KRclunk
	KTremove, KRremove
	KTstat, KRstat
	KTwstat, KRwstat
)
const NOTAG = ^uint16(0)

func str(in string) s {
	return s{
		n:    uint16(len(in)),
		data: []byte(in),
	}
}

//wire9 qid data[13]
//wire9 s n[2] data[n]
//wire9 Hdr size[4] msg[1]

//wire9 Tversion size[4] msg[1] tag[2] msize[4] version[,s]
//wire9 Rversion size[4] msg[1] tag[2] msize[4] version[,s]

var (
	ErrNoConn     = errors.New("no connection")
	ErrBadVersion = errors.New("bad version")
)

type State byte

const (
	StClosed State = iota
	StSyncer
	StSyncee
	StEstablished
	StError
)

type Conn struct {
	txout   chan msg
	netconn net.Conn
	xxx     *bufio.ReadWriter
	err     error
	version string
	state   State
}

func (c *Conn) Read(p []byte) (n int, err error) {
	defer func() { logf("called conn.Read, result n=%d, err=%s", n, err) }()
	return c.xxx.Read(p)
}

func (c *Conn) Write(p []byte) (n int, err error) {
	defer func() { logf("called conn.Write, result n=%d, err=%s", n, err) }()
	return c.xxx.Write(p)
}
func (c *Conn) Flush() (err error) {
	defer func() { logf("called conn.Flush, result err=%s", err) }()
	return c.xxx.Flush()
}

func (c *Conn) Close() (err error) {
	defer func() { logf("called conn.Close, result err=%s", err) }()
	return c.netconn.Close()
}

func Accept(fd net.Listener) (c *Conn, err error) {
	conn, err := fd.Accept()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	c = NewConn(conn)
	return c, negotiateServer(c, &Tversion{
		msg:     byte(KTversion),
		tag:     NOTAG,
		msize:   MaxMsg,
		version: str(Version),
	})
}

func Dial(netw string, addr string) (c *Conn, err error) {
	conn, err := net.Dial(netw, addr)
	if err != nil {
		return nil, err
	}
	c = NewConn(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return c, nil
}

// NewConn opens a new 9p connection from an existing
// conn.
func NewConn(conn net.Conn) (c *Conn) {
	defer func() { go c.run() }()
	return &Conn{
		txout:   make(chan msg),
		netconn: conn,
		xxx: bufio.NewReadWriter(
			bufio.NewReaderSize(conn, MaxMsg),
			bufio.NewWriterSize(conn, MaxMsg),
		),
	}
}

//wire9 Tauth size[4] msg[1] tag[2] afid[4] uname[,s] aname[,s]
//wire9 Rauth size[4] msg[1] tag[2] aqid[13]
//wire9 Rerror size[4] msg[1] tag[2] ename[,s]
//wire9 Tflush size[4] msg[1] tag[2] oldtag[2]
//wire9 Rflush size[4] msg[1] tag[2]
//wire9 Tattach size[4] msg[1] tag[2] fid[4] afid[4] uname[,s] aname[,s]
//wire9 Rattach size[4] msg[1] tag[2] qid[13]
//wire9 Twalk size[4] msg[1] tag[2] fid[4] newfid[4] nwname[2] wname[nwname, []s]
//wire9 Rwalk size[4] msg[1] tag[2] nwqid[2] wqid[nwqid, []qid]
//wire9 Topen size[4] msg[1] tag[2] fid[4] mode[1]
//wire9 Ropen size[4] msg[1] tag[2] qid[13] iounit[4]
//wire9 Topenfd size[4] msg[1] tag[2] fid[4] mode[1]
//wire9 Ropenfd size[4] msg[1] tag[2] qid[13] iounit[4] unixfd[4]
//wire9 Tcreate size[4] msg[1] tag[2] fid[4] name[,s] perm[4] mode[1]
//wire9 Rcreate size[4] msg[1] tag[2] qid[13] iounit[4]
//wire9 Tread size[4] msg[1] tag[2] fid[4] offset[8] count[4]
//wire9 Rread size[4] msg[1] tag[2] count[4] data[count]
//wire9 Twrite size[4] msg[1] tag[2] fid[4] offset[8] count[4] data[count]
//wire9 Rwrite size[4] msg[1] tag[2] count[4]
//wire9 Tclunk size[4] msg[1] tag[2] fid[4]
//wire9 Rclunk size[4] msg[1] tag[2]
//wire9 Tremove size[4] msg[1] tag[2] fid[4]
//wire9 Rremove size[4] msg[1] tag[2]
//wire9 Tstat size[4] msg[1] tag[2] fid[4]
//wire9 Rstat size[4] msg[1] tag[2] stat[size-4-1-2]
//wire9 Twstat size[4] msg[1] tag[2] fid[4] stat[size-4-1-2]
//wire9 Rwstat size[4] msg[1] tag[2]
