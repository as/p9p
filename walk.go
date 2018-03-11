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
