package p9p

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
)

type Kind byte

const (
	KTversion, KRversion byte = iota + 100, iota + 101
	KTauth, KRauth
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

const (
	Version = "9p2000"
	MaxMsg  = 65536
)

func NewBio(conn net.Conn) *Bio {
	return &Bio{
		ReadWriter: bufio.NewReadWriter(
			bufio.NewReaderSize(conn, MaxMsg),
			bufio.NewWriterSize(conn, MaxMsg),
		),
	}
}

var (
	ErrNoConn = errors.New("no connection")
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

type Bio struct {
	*bufio.ReadWriter
	rwc     io.ReadWriteCloser
	err     error
	version string
	state   State
}

func (b *Bio) Version() (string, error) {
	if b.state != StEstablished {
		return "", ErrNoConn
	}
	if b.version == ""{
		return "", ErrBadVersion
	}
	return b.version, nil
}

func Accept(fd net.Listener) (*Bio, error) {
	conn, err := fd.Accept()
	if err != nil {
		return nil, err
	}
	bio := NewBio(conn)

	tv := Tversion{}

	if err = tv.ReadBinary(bio); err != nil {
		panic(err)
	}

	log.Printf("accept: got %#v\n", tv)

	rv := Rversion{
		size:    uint32(4 + 1 + 2 + 4 + len(Version)),
		msg:     KTversion,
		tag:     NOTAG,
		msize:   MaxMsg,
		version: str(Version),
	}
	if err := rv.WriteBinary(bio); err != nil {
		log.Printf("open: err: %s\n", err)
	}

	bio.state = StEstablished
	bio.version = string(rv.version.data)
	
	return bio, bio.Flush()
}

func Dial(netw string, addr string) (*Bio, error) {
	conn, err := net.Dial(netw, addr)
	if err != nil {
		return nil, err
	}

	bio, err := Open(conn)
	if err != nil {
		conn.Close()
	}

	return bio, err
}

func Open(conn net.Conn) (*Bio, error) {
	bio := NewBio(conn)

	tv := Tversion{
		size:    uint32(4 + 1 + 2 + 4 + len(Version)),
		msg:     KTversion,
		tag:     NOTAG,
		msize:   MaxMsg,
		version: str(Version),
	}
	log.Printf("open: sending %#v\n", tv)
	if err := tv.WriteBinary(bio); err != nil {
		log.Printf("open: err: %s\n", err)
	}
	err := bio.Flush()
	if err != nil{
		return nil, err
	}
	
	bio.state = StEstablished
	bio.version = string(tv.version.data)

	return bio, err
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
