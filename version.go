package p9p

type Opts struct {
	Msize   uint32
	Version string
	Cancel  chan bool
}

var defaultOpts = Opts{
	Msize:   65535,
	Version: "9p2000",
	Cancel:  make(chan bool),
}
