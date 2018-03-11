package p9p

type Qids struct {
	List []Qid
}

type Qid struct {
	Type byte
	Ver  uint32
	Path uint64
}
