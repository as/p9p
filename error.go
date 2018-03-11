package p9p

func (c *Conn) Error(why string) (err error) {
	defer func() {
		logf("Error(c->s) (tx=%v): %q %s", why, err)
	}()
	m := &Msg{src: c}
	if !m.writeHeader(KTerror) || !m.writestring(why) {
		return m.err
	}
	c.schedule(m)
	return m.err
}
