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
	logf("schedule: %#v\n", m)
	defer func() { logf("de schedule: %#v\n", m) }()
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
	go func() {
		for {
			m := Msg{src: c}
			if !m.readMsg(0) {
				logf("readmsg: %s\n", m.err)
				panic("TODO(as): handle readMsg failure")
			}
			logf("readmsg loop: %#v\n", m)
			incoming <- m
		}
	}()
	var err error
	for {
		select {
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
			logf("txoutdone: %#v\n", m.Msg)
			inflight[m.Tag] = m.reply
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
