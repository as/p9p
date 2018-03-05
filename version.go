package p9p

import "errors"

var (
	ErrVersion  = errors.New("mismatched version")
	ErrBadMsgRx = errors.New("bad msg reply")
)

/*

type Opts struct {
	Msize   uint32
	Version string
	Cancel  chan bool
}

var defaultOpts = Opts{
	Msize:   65535,
	Version: "9p2000",
	Cancel:  make(chan bool),
}

func (c *Conn) Version(tag int, msize int, v string) (*Rversion, error){
	tx := &Tversion{
	}
	tx.tag = uint16(tag)
	tx.msize = uint32(msize)
	tx.version.n = uint16(len(v))
	tx.version.data = []byte(v)
	tx.size = uint32(4+1+2+4+tx.version.n)

	err := tx.WriteBinary(c.fd)
	if err != nil{
		return nil, err
	}

	rx := &Rversion{}
	err = rx.ReadBinary(c.fd)
	if err != nil{
		return nil, err
	}
	if rx.msg != tx.msg{
		return nil, ErrBadMsgRx
	}
//	if rx.version != tx.version{
//		return nil, ErrVersion
//	}
	panic("not finished")
}
*/
