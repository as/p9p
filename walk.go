package p9p

func (c *Conn) Walk(fid, newfid int, names ...string) (q []Qid, err error) {
	defer func() { logf("Walk: %v %v %v", fid, newfid, names) }()
	m := &Msg{src: c}
	if !m.writeHeader(KTwalk) || !m.writebinary(uint32(fid)) || !m.writebinary(uint32(newfid)) || !m.writeNames(names...) {
		return q, m.err
	}

	if !c.schedule(m) {
		return q, m.err
	}

	m.readQuids(&q)
	return q, m.err
}

// size[4] Tread tag[2] fid[4] offset[8] count[4]
func (c *Conn) ReadFid(fid int, offset int64, p []byte) (n int, err error) {
	var nn int32
	defer func() { logf("Read: %v %v %v", fid, offset, p) }()
	m := &Msg{src: c}
	if !m.writeHeader(KTread) || !m.writebinary(&struct {
		Fid int32
		Ofs int64
		N   int32
	}{int32(fid), offset, int32(len(p))}) {
		return n, m.err
	}

	if !c.schedule(m) || !m.readbinary(&nn) || !m.readbytes(p) {
		return n, m.err
	}

	return int(nn), m.err
}
