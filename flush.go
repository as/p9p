package p9p

func (c *Conn) Flush(oldtag uint) (err error) {
	defer logf("Flush: %s", err)
	m := &Msg{src: c}
	if !m.writeHeader(KTflush) || !m.writebinary(uint16(oldtag)) {
		return m.err
	}

	if !c.schedule(m) {
		return m.err
	}

	x := uint16(0)
	m.readbinary(&x)

	return m.err
}
