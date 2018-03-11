package p9p

func (c *Conn) Error(why string) error {
	defer logf("Error(c->s): %q", why)
	m := &Msg{src: c}
	if !m.writeHeader(KTerror) && !m.writestring(why) && c.schedule(m) {
	}
	return m.err
}
