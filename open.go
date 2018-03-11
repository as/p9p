package p9p

func (c *Conn) Open(fid int, mode byte) (q Qid, iounit uint32, err error) {
	defer func() { logf("Open: %v %s %s", fid, mode, err) }()
	m := &Msg{src: c}
	if !m.writeHeader(KTopen) || !m.writebinary(uint32(fid)) || !m.writebinary(mode) {
		return q, iounit, m.err
	}

	if !c.schedule(m) {
		return q, iounit, m.err
	}

	m.readbinary(&q)
	m.readbinary(&iounit)
	return q, iounit, m.err
}

func (c *Conn) Create(fid int, name string, perm uint32, mode byte) (q Qid, iounit uint32, err error) {
	defer func() { logf("Create: %v %s %s", fid, mode, err) }()
	m := &Msg{src: c}
	if !m.writeHeader(KTcreate) || !m.writebinary(uint32(fid)) || !m.writestring(name) || !m.writebinary(perm) || !m.writebinary(mode) {
		return q, iounit, m.err
	}

	if !c.schedule(m) {
		return q, iounit, m.err
	}

	m.readbinary(&q)
	m.readbinary(&iounit)
	return q, iounit, m.err
}
