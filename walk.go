package p9p

import "io"

func (c *Conn) Remove(fid Fid) (err error) {
	defer func() { logf("Remove: %v", err) }()
	m := &Msg{src: c}
	if !m.writeHeader(KTremove) || !m.writebinary(fid) || !c.schedule(m) {
		return m.err
	}
	return nil
}

func (c *Conn) Clunk(fid Fid) (err error) {
	defer func() { logf("Clunk: %v", err) }()
	m := &Msg{src: c}
	if !m.writeHeader(KTclunk) || !m.writebinary(fid) || !c.schedule(m) {
		return m.err
	}
	return nil
}

func (c *Conn) Walk(fid, newfid Fid, names ...string) (q []Qid, err error) {
	defer func() { logf("Walk: %v %v %v", fid, newfid, names) }()
	m := &Msg{src: c}
	if !m.writeHeader(KTwalk) || !m.writebinary(fid) || !m.writebinary(newfid) || !m.writeNames(names...) {
		return q, m.err
	}

	if !c.schedule(m) {
		return q, m.err
	}

	m.readQuids(&q)
	return q, m.err
}

func (c *Conn) ReadFid(fid Fid, offset int64, p []byte) (n int, err error) {
	var nn int32
	defer func() { logf("Read: %v %v %v", fid, offset, p) }()
	m := &Msg{src: c}
	if !m.writeHeader(KTread) || !m.writebinary(&struct {
		Fid Fid
		Ofs int64
		N   int32
	}{fid, offset, int32(len(p))}) {
		return n, m.err
	}

	if !c.schedule(m) || !m.readbinary(&nn) || !m.readbytes(p) {
		if m.err == io.ErrUnexpectedEOF {
			return n, m.err
		}
	}

	return int(nn), m.err
}

func (c *Conn) WriteFid(fid Fid, offset int64, p []byte) (n int, err error) {
	var nn int32
	defer func() { logf("Write: %v %v %v", fid, offset, p) }()
	m := &Msg{src: c}
	if !m.writeHeader(KTwrite) || !m.writebinary(&struct {
		Fid Fid
		Ofs int64
		N   int32
	}{fid, offset, int32(len(p))}) {
		return n, m.err
	}

	if !m.write(p) || !c.schedule(m) || !m.readbinary(&nn) {
		return n, m.err
	}

	return int(nn), m.err
}
