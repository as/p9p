package p9p

type entry struct {
	tag    uint16
	expect Kind
}

type msg struct {
	*Msg
	reply chan *Msg
}

func (c *Conn) schedule(m *Msg) (ok bool) {
	logf("schedule: %s\n", m)
	defer func() { logf("de schedule: %s\n", m) }()
	if m.err != nil {
		return false
	}

	reply := make(chan *Msg)
	c.txout <- msg{m, reply}
	m0 := <-reply

	if m0.err == nil {
		m0.self()
		*m = *m0
	}
	return m0.err == nil
}

func (c *Conn) run() {
	inflight := make(map[uint16]chan *Msg)
	tag := uint16(0xffff)
	incoming := make(chan Msg)
	fail := make(chan error)
	defer close(fail)
	defer close(c.txout)
	defer c.close(<-c.exit)
	go func() {
		defer close(incoming)
		for {
			select {
			default:
				m := Msg{src: c}
				if !m.readMsg(0) {
					logf("readmsg: %s\n", m.err)
					fail <- m.err
					continue
				}
				logf("readmsg loop: %s\n", m)
				incoming <- m
			case <-c.done:
				return
			}
		}
	}()
	var err error
	for {
		select {
		case <-c.done:
			return
		case err := <-fail:
			if err != nil {
				logf("conn: run: got fatal error: %s", err)
			}
			c.close(<-c.exit)
		case m := <-c.txout:
			if m.Kind == KTversion {
				// According to the 9p man pages, a version message
				// clunks all outstanding fids in flight
				inflight = make(map[uint16]chan *Msg)
				tag = 0xffff
			}
			m.Header.Tag = tag
			err = c.Transmit(m.Msg)
			if err != nil {
				panic("TODO(as): handle transmit failure")
			}
			if m.Header.Kind != KTerror {
				// If the client transmits an error (if that's even possible)
				// there is no need to wait for a reply from the server
				logf("txoutdone: %sn", m.Msg)
				inflight[m.Tag] = m.reply
			}
			tag++
		case m := <-incoming:
			repl, ok := inflight[m.Tag]
			if !ok {
				logf("-incoming: %s\n", m.err)
				panic("TODO(as): handle invalid tag sent as reply")
			}
			repl <- &m
		}
	}
}
