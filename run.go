package p9p

type entry struct {
	tag    uint16
	expect Kind
}

type msg struct {
	*Msg
	reply chan *Msg
}

func (c *Conn) run() {
	inflight := make(map[uint16]chan *Msg)
	ctr := uint16(1)
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
	for {
		select {
		case m := <-c.txout:
			if err := m.Transmit(c); err != nil {
				panic("TODO(as): handle transmit failure")
			}
			logf("txoutdone: %#v\n", m.Msg)
			inflight[m.Tag] = m.reply
			ctr++
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
